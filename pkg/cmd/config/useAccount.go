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
package config

import (
	"errors"
	"github.com/chaosnative/chaosctl/pkg/config"
	"github.com/chaosnative/chaosctl/pkg/types"
	"github.com/chaosnative/chaosctl/pkg/utils"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"os"
)

// useAccountCmd represents the useAccount command
var useAccountCmd = &cobra.Command{
	Use:   "use-account",
	Short: "Sets the current-account and current-username in a chaosconfig file",
	Long:  `Sets the current-account and current-username in a chaosconfig file`,
	Run: func(cmd *cobra.Command, args []string) {
		configFilePath := utils.GetLitmusConfigPath(cmd)

		endpoint, err := cmd.Flags().GetString("endpoint")
		utils.PrintError(err)

		if endpoint == "" {
			prompt := promptui.Select{
				Label: "What's the product?",
				Items: []string{"ChaosNative Cloud", "ChaosNative Enterprise"},
			}

			_, result, err := prompt.Run()

			if err != nil {
				utils.Red.Println(err)
				return
			}

			if result == "ChaosNative Cloud" {
				endpoint = utils.ChaosNativeCloudEndpoint
			} else if result == "ChaosNative Enterprise" {
				validate := func(input string) error {
					if utils.IsValidUrl(input) {
						return nil
					} else {
						return errors.New("Not a valid URL")
					}
				}

				templates := &promptui.PromptTemplates{
					Prompt:  "{{ . }} ",
					Valid:   "{{ . | green }} ",
					Invalid: "{{ . | red }} ",
					Success: "{{ . | bold }} ",
				}

				prompt := promptui.Prompt{
					Label:     "ChaosNative Enterprise Endpoint",
					Templates: templates,
					Validate:  validate,
				}

				endpoint, err = prompt.Run()

				if err != nil {
					utils.Red.Println(errors.New("Prompt err:" + err.Error()))
					return
				}

			}

			for endpoint == "" {
				utils.Red.Println("\n⛔ Host URL can't be empty!!")
				os.Exit(1)
			}
		}

		username, err := cmd.Flags().GetString("username")
		utils.PrintError(err)

		if username == "" {
			prompt := promptui.Prompt{
				Label: "What's the AccessID?",
			}

			username, err = prompt.Run()
			if err != nil {
				utils.Red.Println(errors.New("Prompt err:" + err.Error()))
				os.Exit(1)
			}

			for username == "" {
				utils.Red.Println("\n⛔ AccessID cannot be empty!!")
				os.Exit(1)
			}
		}

		if username == "" || endpoint == "" {
			utils.Red.Println("endpoint or username is not set")
			os.Exit(1)
		}

		exists := config.FileExists(configFilePath)

		err = config.ConfigSyntaxCheck(configFilePath)
		utils.PrintError(err)

		if exists {
			litmusconfig, err := config.YamltoObject(configFilePath)
			utils.PrintError(err)

			isAccountExist := config.IsAccountExists(litmusconfig, username, endpoint)
			if isAccountExist {
				err = config.UpdateCurrent(types.Current{
					CurrentAccount: endpoint,
					CurrentUser:    username,
				}, configFilePath)
				utils.PrintError(err)
			} else {
				utils.Red.Println("\n⛔ Account not exists")
				os.Exit(1)
			}
		} else {
			utils.Red.Println("\n⛔ File not exists")
			os.Exit(1)
		}
	},
}

func init() {
	ConfigCmd.AddCommand(useAccountCmd)
	useAccountCmd.Flags().StringP("username", "u", "", "Help message for toggle")
	useAccountCmd.Flags().StringP("endpoint", "e", "", "Help message for toggle")
}
