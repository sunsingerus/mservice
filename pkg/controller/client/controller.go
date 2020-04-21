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

package controller_client

import (
	"bytes"
	"context"
	"io"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"

	pb "github.com/sunsingerus/mservice/pkg/api/mservice"
)

// SendFile sends file from client to service and receives response back
func SendFile(client pb.MServiceControlPlaneClient, filename string) (int64, error) {
	if _, err := os.Stat(filename); err != nil {
		return 0, err
	}

	log.Infof("Has file %s", filename)
	f, err := os.Open(filename)
	if err != nil {
		log.Infof("ERROR open file %s", filename)
		return 0, err
	}

	log.Infof("START send file %s", filename)
	metadata := pb.NewMetadata(filepath.Base(filename))
	n, err := exchange(client, metadata, f)
	log.Infof("DONE send file %s size %d err %v", filename, n, err)

	return n, err
}

// SendStdin sends STDIN from client to service and receives response back
func SendStdin(client pb.MServiceControlPlaneClient) (int64, error) {
	n, err := exchange(client, nil, os.Stdin)
	log.Infof("DONE send %s size %d err %v", os.Stdin.Name(), n, err)
	return n, err
}

// exchange sends data from client to service and receives response back
func exchange(client pb.MServiceControlPlaneClient, metadata *pb.Metadata, dataSource io.Reader) (n int64, err error) {
	//ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log.Infof("rpcData()")
	rpcData, err := client.Data(ctx)
	if err != nil {
		log.Fatalf("client.Data() failed %v", err)
		os.Exit(1)
	}
	defer func() {
		// This is hand-made flush() replacement for gRPC
		// It is required in order to flush all outstanding data before
		// context's cancel() is called, which simply discards all outstanding data.
		// On receiving end, when cancel() is the first in the race, stream receives 'cancel' and (sometimes) no data
		// instead of complete set of data and EOF
		// See https://github.com/grpc/grpc-go/issues/1714 for more details
		rpcData.CloseSend()
		rpcData.Recv()
	}()

	// Send to server
	log.Infof("Send to Server")

	stream, err := pb.OpenDataChunkStream(rpcData)
	if err != nil {
		log.Fatalf("OpenDataChunkStream() failed %v", err)
		return 0, err
	}
	stream.Type = uint32(pb.DataChunkType_DATA_CHUNK_DATA)
	stream.Metadata = metadata
	stream.UUID_reference = "123"
	stream.Description = "desc"
	io.Copy(stream, dataSource)
	stream.Close()

	// Receive back
	log.Infof("Receive from Server")

	stream, err = pb.OpenDataChunkStream(rpcData)
	if err != nil {
		log.Fatalf("OpenDataChunkStream() failed %v", err)
		return 0, err
	}
	var buf = &bytes.Buffer{}
	n, err = io.Copy(buf, stream)
	stream.Close()
	log.Infof("Incoming filename: %s", stream.Metadata.GetFilename())
	log.Infof("%s", buf.String())

	return n, err
}
