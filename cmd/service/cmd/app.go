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

package cmd

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	pbHealth "github.com/sunsingerus/mservice/pkg/api/health"
	pbMService "github.com/sunsingerus/mservice/pkg/api/mservice"
	"github.com/sunsingerus/mservice/pkg/controller/service/control_plane"
	"github.com/sunsingerus/mservice/pkg/controller/service/health"
	"github.com/sunsingerus/mservice/pkg/version"
)

// CLI parameter variables
var (
	// versionRequest specifies request for version report
	versionRequest bool

	// configFile specifies path to config file to be used
	configFile string

	// serviceAddr specifies address of service to use
	serviceAddress string

	// port specifies port to listen by gRPC handler
	port int
)

func init() {
	flag.BoolVar(&versionRequest, "version", false, "Display version and exit")
	flag.StringVar(&configFile, "config", "", "Path to config file.")
	flag.StringVar(&serviceAddress, "service-address", ":10000", "The address of service to use in the format host:port, as localhost:10000")
	flag.IntVar(&port, "port", 10000, "The server port")

	flag.Parse()
}

// Run is an entry point of the application
func Run() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetLevel(log.TraceLevel)

	if versionRequest {
		fmt.Printf("%s\n", version.Version)
		os.Exit(0)
	}

	// Set OS signals and termination context
	ctx, cancelFunc := context.WithCancel(context.Background())
	stopChan := make(chan os.Signal, 2)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-stopChan
		cancelFunc()
		<-stopChan
		os.Exit(1)
	}()

	log.Infof("Starting service. Version:%s GitSHA:%s BuiltAt:%s\n", version.Version, version.GitSHA, version.BuiltAt)

	log.Infof("Listening on %s", serviceAddress)
	listener, err := net.Listen("tcp", serviceAddress)
	if err != nil {
		log.Fatalf("failed to Listen() %v", err)
		os.Exit(1)
	}

	grpcServer := grpc.NewServer(getGRPCServerOptions()...)
	pbMService.RegisterMServiceControlPlaneServer(grpcServer, &control_plane.Server{})
	pbHealth.RegisterHealthServer(grpcServer, &health.Server{})

	go func() {
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatalf("failed to Serve() %v", err)
			os.Exit(1)
		}
	}()

	log.Infof("Press Ctrl+C to exit...")
	<-ctx.Done()
}

// getGRPCServerOptions builds gRPC server options from flags
func getGRPCServerOptions() []grpc.ServerOption {
	var opts []grpc.ServerOption
	return opts
}
