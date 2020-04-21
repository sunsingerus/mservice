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
	"fmt"
	"os"

	"github.com/MakeNowJust/heredoc"
	"github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	cmd "github.com/spf13/cobra"
	conf "github.com/spf13/viper"
)

const (
	defaultConfigFileName = ".mservice-client.yaml"
	defaultServiceAddress = "localhost:10000"
)

// CLI parameter variables
var (
	// verbose specifies whether app should be verbose
	verbose bool

	// configFile specifies path to config file to be used
	configFile string

	// serviceAddr specifies address of service to use
	serviceAddress string
)

var (
	client grpcClient

	// rootCmd represents the base command when called without any sub-commands
	rootCmd = &cmd.Command{
		Use:   "mservice-client [COMMAND]",
		Short: "A command-line client for mservice.",
		Long: heredoc.Docf(`
			For setting the address of the form HOST:PORT, you can
			- use the flag --%s=%s
			- or you can set '%s: %s' in config $HOME/%s
			`,
			serviceAddressFlagName,
			defaultServiceAddress,
			serviceAddressFlagName,
			defaultServiceAddress,
			defaultConfigFileName,
		),
		PersistentPreRun: func(cmd *cmd.Command, args []string) {
			client = grpcClient{
				address: conf.GetString(serviceAddressFlagName),
			}
			log.Debugf("using address: %s", client.address)
		},
	}
)

const (
	serviceAddressFlagName = "service-address"
)

func init() {
	cmd.OnInitialize(initConfig)
	log.SetFormatter(&log.TextFormatter{})

	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().StringVar(&configFile, "config", "", fmt.Sprintf("config file (default is $HOME/%s)", defaultConfigFileName))
	rootCmd.PersistentFlags().StringVar(&serviceAddress, serviceAddressFlagName, defaultServiceAddress, fmt.Sprintf("The address of service to use in the format host:port, as %s", defaultServiceAddress))
	if err := conf.BindPFlag(serviceAddressFlagName, rootCmd.PersistentFlags().Lookup(serviceAddressFlagName)); err != nil {
		panic(err)
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if verbose {
		log.SetLevel(log.TraceLevel)
	}

	if configFile == "" {
		// Use config file from home directory
		homedir, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		conf.AddConfigPath(homedir)
		configFile = defaultConfigFileName
	}
	conf.SetConfigFile(configFile)
	conf.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := conf.ReadInConfig(); err == nil {
		log.Debugf("using config file: %v", conf.ConfigFileUsed())
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Error(err)
		os.Exit(1)
	}
}

type grpcClient struct {
	address string
}
