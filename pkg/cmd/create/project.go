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
package create

import (
	"errors"
	"github.com/chaosnative/chaosctl/pkg/apis"
	"github.com/chaosnative/chaosctl/pkg/utils"
	"github.com/manifoldco/promptui"
	"os"

	"github.com/spf13/cobra"
)

// projectCmd represents the project command
var projectCmd = &cobra.Command{
	Use: "project",
	Short: `Create a project
	Example:
	#create a project
	chaosctl create project --name new-proj

	Note: The default location of the config file is $HOME/.chaosconfig, and can be overridden by a --config flag
	`,
	Run: func(cmd *cobra.Command, args []string) {
		credentials, err := utils.GetCredentials(cmd)
		utils.PrintError(err)

		projectName, err := cmd.Flags().GetString("name")
		utils.PrintError(err)

		if projectName == "" {
			prompt := promptui.Prompt{
				Label: "Enter a project name",
			}

			projectName, err = prompt.Run()
			if err != nil {
				utils.Red.Println(errors.New("Prompt err:" + err.Error()))
				os.Exit(1)
			}
		}

		userDetails, err := apis.GetProjectDetails(credentials)
		utils.PrintError(err)

		_, err = apis.CreateProjectRequest(userDetails.Data.ID, projectName, credentials)
		utils.PrintError(err)
	},
}

func init() {
	CreateCmd.AddCommand(projectCmd)
	projectCmd.Flags().String("name", "", "Set the project name to create it")
}
