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
	"github.com/sunsingerus/mservice/pkg/transiever"
	"io"
	"os"

	log "github.com/golang/glog"

	pb "github.com/sunsingerus/mservice/pkg/api/mservice"
)

func Init() {
	transiever.Init()
}

func GetOutgoingQueue() chan *pb.Command {
	return transiever.GetOutgoingQueue()
}

func GetIncomingQueue() chan *pb.Command {
	return transiever.GetIncomingQueue()
}

type MServiceControlPlaneEndpoint struct {
	pb.UnimplementedMServiceControlPlaneServer
}

func (s *MServiceControlPlaneEndpoint) Commands(server pb.MServiceControlPlane_CommandsServer) error {
	log.Info("Commands() called")
	defer log.Info("Commands() exited")

	transiever.CommandsExchangeEndlessLoop(server)
	return nil
}

func (s *MServiceControlPlaneEndpoint) Data(server pb.MServiceControlPlane_DataServer) error {
	log.Info("Data() called")
	defer log.Info("Data() exited")

	stream, _ := pb.OpenIncomingDataChunkStream(server)
	defer stream.Close()
	_, err := io.Copy(os.Stdout, stream)

	log.Infof("Incoming filename: %s", stream.Metadata.GetFilename())
	return err
}

func (s *MServiceControlPlaneEndpoint) Metrics(server pb.MServiceControlPlane_MetricsServer) error {
	log.Info("Metrics() called")
	defer log.Info("Metrics() exited")

	return nil
}
