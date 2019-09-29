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
	"errors"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Modify the configuration file",
	Long:  ``,
}

var setPropertyCmd = &cobra.Command{
	Use:   "set PROPERTY_NAME PROPERTY_VALUE",
	Short: "Sets an individual value in the configuration file",
	Long: `Sets an individual value in the configuration file

PROPERTY_NAME is a dot delimited name where each token represents either an attribute name or a map key.  Map keys may not contain dots. 
	
PROPERTY_VALUE is the new value you wish to set.
	`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires a PROPERTY_NAME argument")
		}
		if len(args) < 2 {
			return errors.New("requires a PROPERTY_VALUE argument")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		value := args[1]
		viper.Set(key, value)
		viper.WriteConfig()
	},
}

var unsetPropertyCmd = &cobra.Command{
	Use:   "unset PROPERTY_NAME",
	Short: "Unsets an individual value in the configuration file",
	Long: `Unsets an individual value in the configuration file
	
PROPERTY _NAME is a dot delimited name where each token represents either an attribute name or a map key.  Map keys may not contain dots.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires a PROPERTY_NAME argument")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		delete(viper.Get(key).(map[string]interface{}), "key")
		viper.WriteConfig()
	},
}

var viewConfigFileCmd = &cobra.Command{
	Use:   "view",
	Short: "Display the .gitlabctl.yaml file.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		file, err := ioutil.ReadFile(viper.ConfigFileUsed())
		if err != nil {
			log.Fatal(err)
		}
		fmt.Print(string(file))
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(setPropertyCmd)
	configCmd.AddCommand(unsetPropertyCmd)
	configCmd.AddCommand(viewConfigFileCmd)
}
