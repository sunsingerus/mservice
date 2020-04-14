// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package mservice

import (
	"fmt"
	"io"

	log "github.com/sirupsen/logrus"
)

type DataChunkSenderReceiver interface {
	Send(*DataChunk) error
	Recv() (*DataChunk, error)
}

// DataChunkStream is a handler to open stream of DataChunk's
// Inspired by os.File handler and is expected to be used in the same context.
// DataChunkStream implements the following interfaces:
//	- io.Writer
//	- io.WriterTo
// 	- io.ReaderFrom
//	- io.Closer
// and thus can be used in any functions, which operate these interfaces, such as io.Copy()
type DataChunkStream struct {
	Owner DataChunkSenderReceiver

	Type           uint32
	Name           string
	Metadata       *Metadata
	Version        uint32
	UUID_reference string
	Description    string
	Offset         uint64
}

// OpenDataChunkStream opens DataChunk's stream
// Inspired by os.OpenFile()
func OpenDataChunkStream(owner DataChunkSenderReceiver) (*DataChunkStream, error) {
	return &DataChunkStream{
		Owner: owner,
	}, nil
}

func (s *DataChunkStream) ensureMetadata() {
	if s.Metadata == nil {
		s.Metadata = new(Metadata)
	}
}

// Implements io.Writer
//
// Write writes len(p) bytes from p to the underlying data stream.
// It returns the number of bytes written from p (0 <= n <= len(p))
// and any error encountered that caused the write to stop early.
// Write must return a non-nil error if it returns n < len(p).
// Write must not modify the slice data, even temporarily.
//
// Implementations must not retain p.
func (s *DataChunkStream) Write(p []byte) (n int, err error) {
	n = len(p)
	log.Infof("before Send()")
	var md *Metadata
	if s.Offset == 0 {
		// First chunk in the stream, it may have some metadata
		md = s.Metadata
	}
	dataChunk := NewDataChunk(md, &s.Offset, false, p)
	err = s.Owner.Send(dataChunk)
	if err == io.EOF {
		log.Infof("Send() received EOF, return from func")
		n = 0
	}
	if err != nil {
		log.Fatalf("failed to Send() %v", err)
		n = 0
	}

	s.Offset += uint64(n)

	log.Infof("after Send()")
	return
}

// Implements io.WriterTo
//
// WriterTo is the interface that wraps the WriteTo method.
//
// WriteTo writes data to w until there's no more data to write or
// when an error occurs. The return value n is the number of bytes
// written. Any error encountered during the write is also returned.
//
// The Copy function uses WriterTo if available.
func (s *DataChunkStream) WriteTo(dst io.Writer) (n int64, err error) {
	n = 0
	for {
		dataChunk, readErr := s.Owner.Recv()

		if dataChunk != nil {
			// We have data chunk received
			filename := "not specified"
			if md := dataChunk.GetMetadata(); md != nil {
				filename = md.GetFilename()
				if filename != "" {
					s.ensureMetadata()
					s.Metadata.SetFilename(filename)
				}
			}
			offset := "not specified"
			if off, ok := dataChunk.GetOffsetWithAvailabilityReport(); ok {
				offset = fmt.Sprintf("%d", off)
			}

			len := len(dataChunk.GetBytes())
			log.Infof("Data.Recv() got msg. filename: '%s', chunk len: %d, chunk offset: %s, last chunk: %v",
				filename,
				len,
				offset,
				dataChunk.GetLast(),
			)
			fmt.Printf("%s\n", string(dataChunk.GetBytes()))

			_, _ = dst.Write(dataChunk.GetBytes())
			n += int64(len)

			if dataChunk.GetLast() {
				return
			}
		}

		if readErr == nil {
			// All went well, ready to receive more data
		} else if readErr == io.EOF {
			// Correct EOF
			log.Infof("Data.Recv() get EOF")
			err = readErr
			return
		} else {
			// Stream broken
			log.Infof("Data.Recv() got err: %v", err)
			err = readErr
			return
		}
	}
}

// Implements io.Reader
//
// Read reads up to len(p) bytes into p. It returns the number of bytes
// read (0 <= n <= len(p)) and any error encountered. Even if Read
// returns n < len(p), it may use all of p as scratch space during the call.
// If some data is available but not len(p) bytes, Read conventionally
// returns what is available instead of waiting for more.
//
// When Read encounters an error or end-of-file condition after
// successfully reading n > 0 bytes, it returns the number of
// bytes read. It may return the (non-nil) error from the same call
// or return the error (and n == 0) from a subsequent call.
// An instance of this general case is that a Reader returning
// a non-zero number of bytes at the end of the input stream may
// return either err == EOF or err == nil. The next Read should
// return 0, EOF.
//
// Callers should always process the n > 0 bytes returned before
// considering the error err. Doing so correctly handles I/O errors
// that happen after reading some bytes and also both of the
// allowed EOF behaviors.
//
// Implementations of Read are discouraged from returning a
// zero byte count with a nil error, except when len(p) == 0.
// Callers should treat a return of 0 and nil as indicating that
// nothing happened; in particular it does not indicate EOF.
//
// Implementations must not retain p.
func (s *DataChunkStream) Read(p []byte) (n int, err error) {
	return
}

// Implements io.ReaderFrom
//
// ReaderFrom is the interface that wraps the ReadFrom method.
//
// ReadFrom reads data from r until EOF or error.
// The return value n is the number of bytes read.
// Any error except io.EOF encountered during the read is also returned.
//
// The Copy function uses ReaderFrom if available.
func (s *DataChunkStream) ReadFrom(src io.Reader) (n int64, err error) {
	n = 0
	p := make([]byte, 1024)
	for {
		read, readErr := src.Read(p)
		n += int64(read)
		if read > 0 {
			_, writeErr := s.Write(p[:read])
			if writeErr != nil {
				err = writeErr
				return
			}
		}
		if readErr == io.EOF {
			readErr = nil
			return
		}
		if readErr != nil {
			return
		}
	}
}

// Implements io.Closer
//
// Closer is the interface that wraps the basic Close method.
//
// The behavior of Close after the first call is undefined.
// Specific implementations may document their own behavior.
func (s *DataChunkStream) Close() error {
	if s.Offset == 0 {
		// No data were send earlier, can simply close
		return nil
	}

	// Some data were send earlier, need to finalize transmission first

	// Send "last" data chunk
	dataChunk := NewDataChunk(nil, nil, true, nil)
	err := s.Owner.Send(dataChunk)
	if err == io.EOF {
		log.Infof("Send() received EOF, return from func")
	}
	if err != nil {
		log.Fatalf("failed to Send() %v", err)
	}
	log.Infof("after Send()")

	return err
}
