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

package client_transport

import (
	log "github.com/golang/glog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/testdata"
	"os"
)

func SetupTransport(caFile, serverHostOverride string) ([]grpc.DialOption, error) {
	if caFile == "" {
		caFile = testdata.Path("ca.pem")
	}

	transportCredentials, err := credentials.NewClientTLSFromFile(caFile, serverHostOverride)
	if err != nil {
		log.Fatalf("failed to create TLS credentials %v", err)
		os.Exit(1)
	}

	log.Infof("enabling TLS with ca=%s", caFile)

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(transportCredentials),
	}

	return opts, nil
}
