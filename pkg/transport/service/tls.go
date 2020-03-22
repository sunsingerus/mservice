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

package service_transport

import (
	log "github.com/golang/glog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/testdata"
	"os"
)

func SetupTransport(tlsCertFile, tlsKeyFile string) ([]grpc.ServerOption, error) {
	if tlsCertFile == "" {
		tlsCertFile = testdata.Path("server1.pem")
	}
	if tlsKeyFile == "" {
		tlsKeyFile = testdata.Path("server1.key")
	}

	// TransportCredentials can be created by two ways
	// 1. Directly from files via NewServerTLSFromFile()
	// 2. Or through intermediate Certificate

	// Create TransportCredentials directly from files
	transportCredentials, err := credentials.NewServerTLSFromFile(tlsCertFile, tlsKeyFile)
	// Create TransportCredentials through intermediate Certificate
	// needs "crypto/tls"
	// cert, err := tls.LoadX509KeyPair(testdata.Path("server1.pem"), testdata.Path("server1.key"))
	// transportCredentials := credentials.NewServerTLSFromCert(&cert)

	if err != nil {
		log.Fatalf("failed to generate credentials %v", err)
		os.Exit(1)
	}

	opts := []grpc.ServerOption{
		// Enable TLS transport for connections
		grpc.Creds(transportCredentials),
	}

	return opts, nil
}
