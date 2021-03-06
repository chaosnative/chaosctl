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
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/manifoldco/promptui"

	"github.com/chaosnative/chaosctl/pkg/apis"
	"github.com/chaosnative/chaosctl/pkg/config"
	"github.com/chaosnative/chaosctl/pkg/types"
	"github.com/chaosnative/chaosctl/pkg/utils"
	"github.com/golang-jwt/jwt"
	"github.com/spf13/cobra"
)

// setAccountCmd represents the setAccount command
var setAccountCmd = &cobra.Command{
	Use: "set-account",
	Short: `Sets an account entry in chaosconfig.
		Examples(s)
		#set a new account
		chaosctl config set-account  --endpoint "" --access_id "" --access_key ""
		`,
	Run: func(cmd *cobra.Command, args []string) {

		configFilePath := utils.GetLitmusConfigPath(cmd)

		var (
			authInput types.AuthInput
			err       error
		)

		authInput.Endpoint, err = cmd.Flags().GetString("endpoint")
		utils.PrintError(err)

		authInput.AccessID, err = cmd.Flags().GetString("access_id")
		utils.PrintError(err)

		authInput.AccessKey, err = cmd.Flags().GetString("access_key")
		utils.PrintError(err)

		if authInput.Endpoint == "" {
			prompt := promptui.Select{
				Label: "What's the product name?",
				Items: []string{"Harness Chaos Engineering Cloud", "Harness Chaos Engineering OnPrem"},
			}

			_, result, err := prompt.Run()

			if err != nil {
				utils.Red.Println(err)
				return
			}

			if result == "Harness Chaos Engineering Cloud" {
				authInput.Endpoint = utils.HarnessChaosEngineeringCloudEndpoint
			} else if result == "Harness Chaos Engineering OnPrem" {
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
					Label:     "Harness Chaos Engineering OnPrem Endpoint",
					Templates: templates,
					Validate:  validate,
				}

				authInput.Endpoint, err = prompt.Run()

				if err != nil {
					utils.Red.Printf("Prompt failed %v\n", err)
					return
				}

			}

			utils.White_B.Print(authInput.Endpoint)
			for authInput.Endpoint == "" {
				utils.Red.Println("\n⛔ Host URL can't be empty!!")
				os.Exit(1)
			}

			ep := strings.TrimRight(authInput.Endpoint, "/")
			newUrl, err := url.Parse(ep)
			utils.PrintError(err)

			authInput.Endpoint = newUrl.String()
		}

		if authInput.AccessID == "" {
			prompt := promptui.Prompt{
				Label: "What's the AccessID?",
			}

			authInput.AccessID, err = prompt.Run()
			if err != nil {
				utils.Red.Println(err)
				os.Exit(1)
			}

			if authInput.AccessID == "" {
				utils.Red.Println("\n⛔ AccessID cannot be empty!")
				return
			}
		}

		if authInput.AccessKey == "" {
			prompt := promptui.Prompt{
				Label: "What's the AccessKey?",
				Mask:  '*',
			}

			authInput.AccessKey, err = prompt.Run()
			if err != nil {
				utils.Red.Println(errors.New("Prompt err:" + err.Error()))
				return
			}

			if authInput.AccessKey == "" {
				utils.Red.Println("\n⛔ AccessKey cannot be empty!")
				return
			}
		}

		if authInput.Endpoint != "" && authInput.AccessID != "" && authInput.AccessKey != "" {
			exists := config.FileExists(configFilePath)
			var lgt int
			if exists {
				lgt, err = config.GetFileLength(configFilePath)
				utils.PrintError(err)
			}

			resp, err := apis.Auth(authInput)
			utils.PrintError(err)
			// Decoding token
			token, _ := jwt.Parse(resp.AccessToken, nil)
			if token == nil {
				os.Exit(1)
			}
			claims, _ := token.Claims.(jwt.MapClaims)

			var user = types.User{
				ExpiresIn: fmt.Sprint(time.Now().Add(time.Second * time.Duration(resp.ExpiresIn)).Unix()),
				Token:     resp.AccessToken,
				Username:  claims["username"].(string),
			}

			var users []types.User
			users = append(users, user)

			var account = types.Account{
				Endpoint: authInput.Endpoint,
				Users:    users,
			}

			// If config file doesn't exist or length of the file is zero.
			if !exists || lgt == 0 {

				var accounts []types.Account
				accounts = append(accounts, account)

				var litmuCtlConfig = types.LitmuCtlConfig{
					APIVersion:     "v1",
					Kind:           "Config",
					CurrentAccount: authInput.Endpoint,
					CurrentUser:    claims["username"].(string),
					Accounts:       accounts,
				}

				err := config.CreateNewLitmusCtlConfig(configFilePath, litmuCtlConfig)
				utils.PrintError(err)

				os.Exit(0)
			} else {
				// checking syntax
				err = config.ConfigSyntaxCheck(configFilePath)
				utils.PrintError(err)

				var updateLitmusCtlConfig = types.UpdateLitmusCtlConfig{
					Account:        account,
					CurrentAccount: authInput.Endpoint,
					CurrentUser:    claims["username"].(string),
				}

				err = config.UpdateLitmusCtlConfig(updateLitmusCtlConfig, configFilePath)
				utils.PrintError(err)
			}
			utils.White_B.Printf("\naccount.accessID/%s configured\n", claims["username"].(string))

		} else {
			utils.Red.Println("\nError: some flags are missing. Run 'chaosctl config set-account --help' for usage. ")
		}
	},
}

func init() {
	ConfigCmd.AddCommand(setAccountCmd)

	setAccountCmd.Flags().StringP("endpoint", "e", "", "Account endpoint. Mandatory")
	setAccountCmd.Flags().String("access_id", "", "User AccessID. Mandatory")
	setAccountCmd.Flags().String("access_key", "", "User AccessKey. Mandatory")
}
