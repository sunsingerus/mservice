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

// DataChunkTransporter defines transport level interface (for both client and server),
// which serves DataChunk streams bi-directionally.
type DataChunkTransporter interface {
	Send(*DataChunk) error
	Recv() (*DataChunk, error)
}

// DataChunkFile is a handler to set of DataChunk(s)
// Inspired by os.File handler and is expected to be used in the same context.
// DataChunkFile implements the following interfaces:
//	- io.Writer
//	- io.WriterTo
// 	- io.ReaderFrom
//	- io.Closer
// and thus can be used in any functions, which operate these interfaces, such as io.Copy()
type DataChunkFile struct {
	transportServer DataChunkTransporter

	Type          uint32
	Name          string
	Metadata      *Metadata
	Version       uint32
	UUIDReference string
	Description   string
	Offset        uint64
}

// OpenDataChunkFile opens set of DataChunk(s)
// Inspired by os.OpenFile()
func OpenDataChunkFile(transportServer DataChunkTransporter) (*DataChunkFile, error) {
	return &DataChunkFile{
		transportServer: transportServer,
	}, nil
}

// ensureMetadata ensures DataChunkFile has Metadata in place
func (s *DataChunkFile) ensureMetadata() {
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
func (s *DataChunkFile) Write(p []byte) (n int, err error) {
	n = len(p)
	log.Infof("Write() start")
	var md *Metadata = nil
	if s.Offset == 0 {
		// First chunk in the stream, it may have some metadata
		md = s.Metadata
	}
	chunk := NewDataChunk(md, &s.Offset, false, p)
	err = s.transportServer.Send(chunk)
	if err != nil {
		// We have some kind of error and were not able to send the data
		n = 0
		if err == io.EOF {
			// Not sure, probably this is not an error, after all
			log.Infof("Send() received EOF")
		} else {
			log.Fatalf("failed to Send() %v", err)
		}
	}

	s.Offset += uint64(n)

	log.Infof("Write() end")
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
func (s *DataChunkFile) WriteTo(dst io.Writer) (n int64, err error) {
	n = 0
	for {
		// Whether this chunk is the last one on the stream
		lastChunk := false
		dataChunk, readErr := s.transportServer.Recv()

		if dataChunk != nil {
			// We've got data chunk and it has to be processed no matter what.
			// Even in case Recv() reported error, we still need to process obtained data

			// Fetch filename from the chunks stream - it may be in any chunk, actually
			filename := "not specified"
			if md := dataChunk.GetMetadata(); md != nil {
				filename = md.GetFilename()
				if filename != "" {
					s.ensureMetadata()
					s.Metadata.SetFilename(filename)
				}
			}

			// Fetch offset of this chunk within the stream
			offset := "not specified"
			if off, ok := dataChunk.GetOffsetWithAvailabilityReport(); ok {
				offset = fmt.Sprintf("%d", off)
			}

			// How many bytes do we have in this chunk?
			len := len(dataChunk.GetBytes())
			log.Infof("Data.Recv() got msg. filename: '%s', chunk len: %d, chunk offset: %s, last chunk: %v",
				filename,
				len,
				offset,
				dataChunk.GetLast(),
			)
			fmt.Printf("%s\n", string(dataChunk.GetBytes()))

			// TODO need to handle write errors
			_, _ = dst.Write(dataChunk.GetBytes())

			n += int64(len)

			// This is officially last chunk on the chunks stream, so no need to Recv() any more, break the recv loop,
			// but we need to handle possible errors first
			if dataChunk.GetLast() {
				lastChunk = true
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

		if lastChunk {
			// This is officially last chunk on the chunks stream, so no need to Recv() any more, break the recv loop
			break
		}
	}

	return
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
func (s *DataChunkFile) Read(p []byte) (n int, err error) {
	return 0, fmt.Errorf("not implemented function: Read()")
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
func (s *DataChunkFile) ReadFrom(src io.Reader) (n int64, err error) {
	n = 0
	p := make([]byte, 1024)
	for {
		read, readErr := src.Read(p)
		n += int64(read)
		if read > 0 {
			// We've got some data and it has to be processed no matter what
			// Write these bytes to data stream
			_, writeErr := s.Write(p[:read])
			if writeErr != nil {
				err = writeErr
				return
			}
		}

		if readErr != nil {
			// We have some kind of an error
			if readErr == io.EOF {
				// However, EOF is not an error on read
				readErr = nil
			}
			// In case of an error or EOF, break the read loop
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
func (s *DataChunkFile) Close() error {
	if s.Offset == 0 {
		// No data were sent via this stream, no need to send finalizer
		return nil
	}

	// Some data were sent via this stream, need to finalize transmission with finalizer packet

	// Send "last" data chunk
	chunk := NewDataChunk(nil, nil, true, nil)
	err := s.transportServer.Send(chunk)
	if err != nil {
		if err == io.EOF {
			log.Infof("Send() received EOF")
		} else {
			log.Fatalf("failed to Send() %v", err)
		}
	}
	log.Infof("Close() done")

	return err
}

// TODO implement Reset function for DataChunkFile,
//  so the same descriptor can be used for multiple transmissions.
