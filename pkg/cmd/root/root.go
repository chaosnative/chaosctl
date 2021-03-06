/*
Copyright © 2021 The LitmusChaos Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package rootCmd

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/chaosnative/chaosctl/pkg/cmd/connect"
	"github.com/chaosnative/chaosctl/pkg/cmd/hubgen"

	"github.com/chaosnative/chaosctl/pkg/cmd/upgrade"
	"github.com/chaosnative/chaosctl/pkg/cmd/version"
	config2 "github.com/chaosnative/chaosctl/pkg/config"
	"github.com/chaosnative/chaosctl/pkg/utils"

	"github.com/chaosnative/chaosctl/pkg/cmd/config"
	"github.com/chaosnative/chaosctl/pkg/cmd/create"
	"github.com/chaosnative/chaosctl/pkg/cmd/get"
	"github.com/spf13/cobra"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "chaosctl",
	Short: "ChaosCTL controls the ChaosNative cloud chaos delegate plane",
	Long:  `ChaosCTL controls the ChaosNative cloud chaos delegate plane. ` + "\n" + ` Find more information at: https://github.com/chaosnative/chaosctl`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.AddCommand(config.ConfigCmd)
	rootCmd.AddCommand(connect.ConnectCmd)
	rootCmd.AddCommand(create.CreateCmd)
	rootCmd.AddCommand(get.GetCmd)
	rootCmd.AddCommand(version.VersionCmd)
	rootCmd.AddCommand(upgrade.UpgradeCmd)
	rootCmd.AddCommand(hubgen.HubgenCmd)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.chaosconfig)")
	rootCmd.PersistentFlags().BoolVar(&config2.SkipSSLVerify, "skip-ssl", false, "skip-ssl, chaosctl will skip ssl/tls verification while communicating with portal")
	rootCmd.PersistentFlags().StringVar(&config2.CACert, "cacert", "", "cacert <path_to_crt_file> , custom ca certificate used for communicating with portal")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".chaosconfig" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(utils.DefaultFileName)
	}

	viper.AutomaticEnv() // read in environment variables that match

	if config2.SkipSSLVerify {
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	} else if config2.CACert != "" {
		caCert, err := ioutil.ReadFile(config2.CACert)
		cobra.CheckErr(err)
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{RootCAs: caCertPool}
	}

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
