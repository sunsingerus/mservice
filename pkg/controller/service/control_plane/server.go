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

package control_plane

import (
	"bytes"
	"io"
	"strings"

	log "github.com/sirupsen/logrus"
	pb "github.com/sunsingerus/mservice/pkg/api/mservice"
)

// Server is an implementation of gRPC MServiceControlPlane service
type Server struct {
	pb.UnimplementedMServiceControlPlaneServer
}

// Data is a data processing function in gRPC MServiceControlPlane Service
// It received string, makes it uppercase and sends back to client
func (s *Server) Data(server pb.MServiceControlPlane_DataServer) error {
	log.Info("Data() called")
	defer log.Info("Data() exited")

	log.Infof("Receive from Client")

	// Open incoming data chunks stream as file
	f, err := pb.OpenDataChunkFile(server)
	if err != nil {
		log.Fatalf("OpenDataChunkFile() failed %v", err)
		return err
	}
	var buf = &bytes.Buffer{}
	_, err = io.Copy(buf, f)
	log.Infof("Incoming filename: %s", f.Metadata.GetFilename())
	log.Infof("%s", buf.String())
	f.Close()

	// Process incoming data
	buf = s.process(buf)

	// Send response
	log.Infof("Send to Client")
	f, err = pb.OpenDataChunkFile(server)
	if err != nil {
		log.Fatalf("OpenDataChunkFile() failed %v", err)
		return err
	}
	f.Type = uint32(pb.DataChunkType_DATA_CHUNK_DATA)
	f.Metadata = pb.NewMetadata("returnback.file")
	f.UUIDReference = "321"
	f.Description = "csed"

	_, err = io.Copy(f, buf)
	f.Close()

	return err
}

// process performs actual data processing
func (s *Server) process(buf *bytes.Buffer) *bytes.Buffer {
	var res = &bytes.Buffer{}
	res.WriteString(strings.ToUpper(buf.String()))
	return res
}
