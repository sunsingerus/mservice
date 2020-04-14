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

package transiever_service

import (
	"bytes"
	"io"
	"strings"

	log "github.com/sirupsen/logrus"

	pb "github.com/sunsingerus/mservice/pkg/api/mservice"
)

type MServiceControlPlaneEndpoint struct {
	pb.UnimplementedMServiceControlPlaneServer
}

func (s *MServiceControlPlaneEndpoint) Data(server pb.MServiceControlPlane_DataServer) error {
	log.Info("Data() called")
	defer log.Info("Data() exited")

	// Receive data

	log.Infof("Receive from Client")

	stream, err := pb.OpenDataChunkStream(server)
	if err != nil {
		log.Fatalf("OpenIncomingDataChunkStream() failed %v", err)
		return err
	}
	var buf = &bytes.Buffer{}
	_, err = io.Copy(buf, stream)
	log.Infof("Incoming filename: %s", stream.Metadata.GetFilename())
	log.Infof("%s", buf.String())
	stream.Close()

	// Send back
	var buf2 = &bytes.Buffer{}
	buf2.WriteString(strings.ToUpper(buf.String()))

	log.Infof("Send to Client")
	stream, err = pb.OpenDataChunkStream(server)
	if err != nil {
		log.Fatalf("OpenIncomingDataChunkStream() failed %v", err)
		return err
	}
	stream.Type = uint32(pb.DataChunkType_DATA_CHUNK_DATA)
	stream.Metadata = pb.NewMetadata("returnback.file")
	stream.UUID_reference = "321"
	stream.Description = "csed"

	io.Copy(stream, buf2)
	stream.Close()

	return err
}
