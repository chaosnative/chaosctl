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
package get

import (
	"os"
	"text/tabwriter"

	"github.com/chaosnative/chaosctl/pkg/apis"
	"github.com/chaosnative/chaosctl/pkg/utils"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

// agentsCmd represents the agents command
var agentsCmd = &cobra.Command{
	Use:   "agents",
	Short: "Display list of chaos delegates within the project",
	Long:  `Display list of chaos delegates within the project`,
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
				utils.Red.Println(err)
				os.Exit(1)
			}

			for projectID == "" {
				utils.Red.Println("⛔ Project ID can't be empty!!")
				os.Exit(1)
			}
		}

		agents, err := apis.GetAgentList(credentials, projectID)
		utils.PrintError(err)

		output, err := cmd.Flags().GetString("output")
		utils.PrintError(err)

		switch output {
		case "json":
			utils.PrintInJsonFormat(agents.Data)

		case "yaml":
			utils.PrintInYamlFormat(agents.Data)

		case "":

			writer := tabwriter.NewWriter(os.Stdout, 4, 8, 1, '\t', 0)
			utils.White_B.Fprintln(writer, "AGENTID\tAGENTNAME\tSTATUS\t")

			for _, agent := range agents.Data.GetAgent {
				var status string
				if agent.IsActive {
					status = "ACTIVE"
				} else {
					status = "INACTIVE"
				}
				utils.White.Fprintln(writer, agent.ClusterID+"\t"+agent.AgentName+"\t"+status+"\t")
			}
			writer.Flush()
		}
	},
}

func init() {
	GetCmd.AddCommand(agentsCmd)

	agentsCmd.Flags().String("project-id", "", "Set the project-id. To retrieve projects. Apply `chaosctl get projects`")

	agentsCmd.Flags().StringP("output", "o", "", "Output format. One of:\njson|yaml")
}
