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
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"
	cmd "github.com/spf13/cobra"
	conf "github.com/spf13/viper"
	"google.golang.org/grpc"

	pb "github.com/sunsingerus/mservice/pkg/api/mservice"
	controller "github.com/sunsingerus/mservice/pkg/controller/client"
	"github.com/sunsingerus/mservice/pkg/version"
)

const (
	sendFileFlagName  = "file"
	sendSTDINFlagName = "stdin"
)

var (
	// readFilename specifies file to read and send to service
	sendFilename string

	// readSTDIN specifies whether to read STDIN
	sendSTDIN bool
)

var sendCmd = &cmd.Command{
	Use:   "send",
	Short: "Send file or STDIN to service",
	Args: func(cmd *cmd.Command, args []string) error {
		//if len(args) < 1 {
		//	return errors.New("requires an filename as argument")
		//}
		return nil
	},
	Run: func(cmd *cmd.Command, args []string) {
		//filename := args[0]

		// Set OS signals and termination context
		_, cancelFunc := context.WithCancel(context.Background())
		stopChan := make(chan os.Signal, 2)
		signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)
		go func() {
			<-stopChan
			cancelFunc()
			<-stopChan
			os.Exit(1)
		}()

		log.Infof("Starting client. Version:%s GitSHA:%s BuiltAt:%s\n", version.Version, version.GitSHA, version.BuiltAt)

		log.Infof("Dial() to %s", serviceAddress)
		conn, err := grpc.Dial(serviceAddress, getDialOptions()...)
		if err != nil {
			log.Fatalf("fail to dial %v", err)
			os.Exit(1)
		}
		defer conn.Close()

		client := pb.NewMServiceControlPlaneClient(conn)

		if sendFilename != "" {
			controller.SendFile(client, sendFilename)
		}

		if sendSTDIN {
			controller.SendStdin(client)
		}
	},
}

func init() {
	sendCmd.PersistentFlags().StringVar(&sendFilename, sendFileFlagName, "", "Send file")
	if err := conf.BindPFlag(sendFileFlagName, rootCmd.PersistentFlags().Lookup(sendFileFlagName)); err != nil {
		panic(err)
	}
	sendCmd.PersistentFlags().BoolVar(&sendSTDIN, sendSTDINFlagName, false, "Read data from STDIN and send it")
	if err := conf.BindPFlag(sendSTDINFlagName, rootCmd.PersistentFlags().Lookup(sendSTDINFlagName)); err != nil {
		panic(err)
	}

	rootCmd.AddCommand(sendCmd)
}

// getDialOptions builds gRPC dial options
func getDialOptions() []grpc.DialOption {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithBlock())
	opts = append(opts, grpc.WithInsecure())
	return opts
}
