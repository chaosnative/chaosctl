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
package connect

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/chaosnative/chaosctl/pkg/agent"
	"github.com/chaosnative/chaosctl/pkg/apis"
	"github.com/chaosnative/chaosctl/pkg/k8s"
	"github.com/chaosnative/chaosctl/pkg/types"
	"github.com/chaosnative/chaosctl/pkg/utils"

	"github.com/spf13/cobra"
)

// agentCmd represents the chaos delegate command
var agentCmd = &cobra.Command{
	Use: "chaos-delegate",
	Short: `connect an external chaos delegate.
	Example(s):
	#connect a chaos delegate
	chaosctl connect chaos-delegate --name="new-chaos-delegate" --non-interactive

	#connect a chaos delegate within a project
	chaosctl connect chaos-delegate --name="new-chaos-delegate" --project-id="d861b650-1549-4574-b2ba-ab754058dd04" --non-interactive
	
	Note: The default location of the config file is $HOME/.chaosconfig, and can be overridden by a --config flag
`,
	Run: func(cmd *cobra.Command, args []string) {
		credentials, err := utils.GetCredentials(cmd)
		utils.PrintError(err)

		nonInteractive, err := cmd.Flags().GetBool("non-interactive")
		utils.PrintError(err)

		kubeconfig, err := cmd.Flags().GetString("kubeconfig")
		utils.PrintError(err)

		var newAgent types.Agent

		newAgent.ProjectId, err = cmd.Flags().GetString("project-id")
		utils.PrintError(err)

		userDetails, err := apis.GetProjectDetails(credentials)
		utils.PrintError(err)

		// If projectID is not passed, then it creates a random project
		if newAgent.ProjectId == "" {
			var projectExists = false
		outerloop:
			for _, project := range userDetails.Data.Projects {
				for _, member := range project.Members {
					if (member.UserID == userDetails.Data.ID) && (member.Role == "Owner" || member.Role == "Editor" || member.Role == "Admin") {
						projectExists = true
						break outerloop
					}
				}
			}

			if !projectExists {
				utils.White_B.Print("Creating a random project...")
				newAgent.ProjectId = agent.CreateRandomProject(userDetails.Data.ID, credentials)
			}
		}

		if nonInteractive {

			newAgent.Mode, err = cmd.Flags().GetString("installation-mode")
			utils.PrintError(err)

			if newAgent.Mode == "" {
				utils.Red.Print("Error: --installation-mode flag is empty")
				os.Exit(1)
			}

			newAgent.AgentName, err = cmd.Flags().GetString("name")
			utils.PrintError(err)

			if newAgent.AgentName == "" {
				utils.Red.Print("Error: --name flag is empty")
				os.Exit(1)
			}

			newAgent.Description, err = cmd.Flags().GetString("description")
			utils.PrintError(err)

			newAgent.SkipSSL, err = cmd.Flags().GetBool("skip-ssl")
			utils.PrintError(err)

			newAgent.PlatformName, err = cmd.Flags().GetString("platform-name")
			utils.PrintError(err)

			if newAgent.PlatformName == "" {
				utils.Red.Print("Error: --platform-name flag is empty")
				os.Exit(1)
			}

			newAgent.ClusterType, err = cmd.Flags().GetString("chaos-delegate-type")
			utils.PrintError(err)
			if newAgent.ClusterType == "" {
				utils.Red.Print("Error: --chaos-delegate-type flag is empty")
				os.Exit(1)
			}

			newAgent.NodeSelector, err = cmd.Flags().GetString("node-selector")
			utils.PrintError(err)
			if newAgent.NodeSelector != "" {
				if ok := utils.CheckKeyValueFormat(newAgent.NodeSelector); !ok {
					os.Exit(1)
				}
			}

			toleration, err := cmd.Flags().GetString("tolerations")
			utils.PrintError(err)

			if toleration != "" {
				var tolerations []types.Toleration
				err := json.Unmarshal([]byte(toleration), &tolerations)
				utils.PrintError(err)

				str := "["
				for _, tol := range tolerations {
					str += "{"
					if tol.TolerationSeconds > 0 {
						str += "tolerationSeconds: " + fmt.Sprint(tol.TolerationSeconds) + " "
					}
					if tol.Effect != "" {
						str += "effect: \\\"" + tol.Effect + "\\\" "
					}
					if tol.Key != "" {
						str += "key: \\\"" + tol.Key + "\\\" "
					}

					if tol.Value != "" {
						str += "value: \\\"" + tol.Value + "\\\" "
					}

					if tol.Operator != "" {
						str += "operator : \\\"" + tol.Operator + "\\\" "
					}

					str += " }"
				}
				str += "]"

				newAgent.Tolerations = str
			}

			newAgent.Namespace, err = cmd.Flags().GetString("namespace")
			utils.PrintError(err)

			newAgent.ServiceAccount, err = cmd.Flags().GetString("service-account")
			utils.PrintError(err)

			newAgent.NsExists, err = cmd.Flags().GetBool("ns-exists")
			utils.PrintError(err)

			newAgent.SAExists, err = cmd.Flags().GetBool("sa-exists")
			utils.PrintError(err)

			if newAgent.Mode == "" {
				newAgent.Mode = utils.DefaultMode
			}

			if newAgent.ProjectId == "" {
				utils.Red.Println("Error: --project-id flag is empty")
				os.Exit(1)
			}

			// Check if user has sufficient permissions based on mode
			utils.White_B.Print("\n🏃 Running prerequisites check....")
			agent.ValidateSAPermissions(newAgent.Namespace, newAgent.Mode, &kubeconfig)

			agents, err := apis.GetAgentList(credentials, newAgent.ProjectId)
			utils.PrintError(err)

			// Duplicate agent check
			var isAgentExist = false
			for i := range agents.Data.GetAgent {
				if newAgent.AgentName == agents.Data.GetAgent[i].AgentName {
					utils.White_B.Print(agents.Data.GetAgent[i].AgentName)
					isAgentExist = true
				}
			}

			if isAgentExist {
				utils.Red.Print("Chaos delegate name already exist")
				os.Exit(1)
			}

		} else {

			if newAgent.ProjectId == "" {
				newAgent.ProjectId = agent.GetProjectID(userDetails)
			}

			modeType := agent.GetModeType()

			// Check if user has sufficient permissions based on mode
			utils.White_B.Print("\n🏃 Running prerequisites check....")
			agent.ValidateSAPermissions(newAgent.Namespace, modeType, &kubeconfig)
			newAgent, err = agent.GetAgentDetails(modeType, newAgent.ProjectId, credentials, &kubeconfig)
			utils.PrintError(err)

			newAgent.ServiceAccount, newAgent.SAExists = k8s.ValidSA(newAgent.Namespace, &kubeconfig)
			newAgent.Mode = modeType
		}

		agent.Summary(newAgent, &kubeconfig)

		if !nonInteractive {
			agent.ConfirmInstallation()
		}
		agent, err := apis.ConnectAgent(newAgent, credentials)
		if err != nil {
			utils.Red.Println("\n❌ Chaos delegate connection failed: " + err.Error() + "\n")
			os.Exit(1)
		}
		if agent.Data.UserAgentReg.Token != "" {
			path := fmt.Sprintf("%s/%s/%s.yaml", credentials.Endpoint, utils.ChaosYamlPath, agent.Data.UserAgentReg.Token)
			utils.White_B.Print("Applying YAML:\n", path)
		} else {
			utils.Red.Print("\n🚫 Token Generation failed, chaos delegate installation failed\n")
			os.Exit(1)
		}

		//Apply agent connection yaml
		yamlOutput, err := k8s.ApplyYaml(k8s.ApplyYamlPrams{
			Token:    agent.Data.UserAgentReg.Token,
			Endpoint: credentials.Endpoint,
			YamlPath: utils.ChaosYamlPath,
		}, kubeconfig, false)

		if err != nil {
			utils.Red.Print("\n❌ Failed in applying connection yaml: \n" + err.Error() + "\n")
			os.Exit(1)
		}

		utils.White_B.Print("\n", yamlOutput)

		// Watch subscriber pod status
		k8s.WatchPod(k8s.WatchPodParams{Namespace: newAgent.Namespace, Label: utils.ChaosAgentLabel}, &kubeconfig)

		utils.White_B.Println("\n🚀 Chaos Delegate Connection Successful!! 🎉")
		utils.White_B.Println("👉 Chaos Delegates can be accessed here: " + fmt.Sprintf("%s/%s", credentials.Endpoint, utils.ChaosAgentPath))
	},
}

func init() {
	ConnectCmd.AddCommand(agentCmd)

	agentCmd.Flags().BoolP("non-interactive", "n", false, "Set it to true for non interactive mode | Note: Always set the boolean flag as --non-interactive=Boolean")
	agentCmd.Flags().StringP("kubeconfig", "k", "", "Set to pass kubeconfig file if it is not in the default location ($HOME/.kube/config)")
	agentCmd.Flags().String("tolerations", "", "Set to pass kubeconfig file if it is not in the default location ($HOME/.kube/config)")

	agentCmd.Flags().String("project-id", "", "Set the project-id to install chaos-delegate for the particular project. To see the projects, apply chaosctl get projects")
	agentCmd.Flags().String("installation-mode", "cluster", "Set the installation mode for the kind of chaos-delegate | Supported=cluster/namespace")
	agentCmd.Flags().String("name", "", "Set the chaos-delegate name")
	agentCmd.Flags().String("description", "---", "Set the chaos-delegate description")
	agentCmd.Flags().Bool("skip-ssl", false, "Set whether agent will skip ssl/tls check (can be used for self-signed certs, if cert is not provided in portal)")
	agentCmd.Flags().String("platform-name", "Others", "Set the platform name. Supported- AWS/GKE/Openshift/Rancher/Others")
	agentCmd.Flags().String("chaos-delegate-type", "external", "Set the chaos-delegate-type to external for external chaos-delegates | Supported=external/internal")
	agentCmd.Flags().String("node-selector", "", "Set the node-selector for chaos-delegate components | Format: \"key1=value1,key2=value2\")")
	agentCmd.Flags().String("namespace", "litmus", "Set the namespace for the chaos-delegate installation")
	agentCmd.Flags().String("service-account", "litmus", "Set the service account to be used by the chaos-delegate")
	agentCmd.Flags().Bool("ns-exists", false, "Set the --ns-exists=false if the namespace mentioned in the --namespace flag is not existed else set it to --ns-exists=true | Note: Always set the boolean flag as --ns-exists=Boolean")
	agentCmd.Flags().Bool("sa-exists", false, "Set the --sa-exists=false if the service-account mentioned in the --service-account flag is not existed else set it to --sa-exists=true | Note: Always set the boolean flag as --sa-exists=Boolean\"\n")
}
