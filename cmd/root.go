/*
Copyright © 2019 The Mosteroid Authors.

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
package cmd

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/mosteroid/gitlab-cli/client"
	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gitlab-cli",
	Short: "Command line interface for gitlab",
	Long:  ``,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if viper.GetBool("gitlab.insecure") {
			http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		}
		client.InitClient(viper.GetString("gitlab.baseUrl"), viper.GetString("gitlab.accessToken"))
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "Set the config file (default is $HOME/.gitlab-cli.yaml)")

	rootCmd.PersistentFlags().BoolP("insecure", "k", false, "Allow connections to SSL sites without certs")
	viper.BindPFlag("gitlab.insecure", rootCmd.PersistentFlags().Lookup("insecure"))

	rootCmd.PersistentFlags().String("baseUrl", "", "Set the gitlab base url")
	viper.BindPFlag("gitlab.baseUrl", rootCmd.PersistentFlags().Lookup("baseUrl"))

	rootCmd.PersistentFlags().String("accessToken", "", "Set the user access token")
	viper.BindPFlag("gitlab.accessToken", rootCmd.PersistentFlags().Lookup("accessToken"))

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".gitlab-cli" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".gitlab-cli")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}
}
