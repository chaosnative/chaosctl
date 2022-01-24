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
package upgrade

import (
	"context"
	"errors"
	"github.com/manifoldco/promptui"
	"os"

	"github.com/chaosnative/chaosctl/pkg/apis"
	"github.com/chaosnative/chaosctl/pkg/utils"
	"github.com/spf13/cobra"
)

// createCmd represents the create command
var agentCmd = &cobra.Command{
	Use:   "agent",
	Short: `Upgrades the ChaosNative Cloud agent plane.`,
	Run: func(cmd *cobra.Command, args []string) {
		credentials, err := utils.GetCredentials(cmd)
		utils.PrintError(err)

		projectID, err := cmd.Flags().GetString("project-id")
		utils.PrintError(err)

		if projectID == "" {
			prompt := promptui.Prompt{
				Label: "What's the ProjectID?",
			}

			projectID, err = prompt.Run()
			if err != nil {
				utils.Red.Println(errors.New("Prompt err:" + err.Error()))
				os.Exit(1)
			}
		}

		cluster_id, err := cmd.Flags().GetString("cluster-id")
		utils.PrintError(err)

		if cluster_id == "" {
			prompt := promptui.Prompt{
				Label: "What's the ClusterID?",
			}

			cluster_id, err = prompt.Run()
			if err != nil {
				utils.Red.Println(errors.New("Prompt err:" + err.Error()))
				os.Exit(1)
			}
		}

		_, err = apis.UpgradeAgent(context.Background(), credentials, projectID, cluster_id)
		utils.PrintError(err)
	},
}

func init() {
	UpgradeCmd.AddCommand(agentCmd)
	agentCmd.Flags().String("project-id", "", "Enter the project ID")
	agentCmd.Flags().String("cluster-id", "", "Enter the cluster ID")
}
