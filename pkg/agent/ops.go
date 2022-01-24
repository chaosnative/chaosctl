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
package agent

import (
	"errors"
	"github.com/chaosnative/chaosctl/pkg/apis"
	"github.com/chaosnative/chaosctl/pkg/k8s"
	"github.com/chaosnative/chaosctl/pkg/types"
	"github.com/chaosnative/chaosctl/pkg/utils"
	"github.com/manifoldco/promptui"
	"os"
	"strconv"
)

// GetProject display list of projects and returns the project id based on input
func GetProjectID(u apis.ProjectDetails) string {
	var projectNames []string
	for _, v := range u.Data.Projects {
		projectNames = append(projectNames, v.Name)
	}

	prompt := promptui.Select{
		Label: "Select a project from the list",
		Items: projectNames,
		Size:  len(projectNames),
	}

	counter, _, err := prompt.Run()
	utils.Red.Println(errors.New("Prompt err:" + err.Error()))

	return u.Data.Projects[counter].ID
}

// GetMode gets mode of agent installation as input
func GetModeType() string {
	prompt := promptui.Select{
		Label: "What's the installation mode?",
		Items: []string{"Cluster", "Namespace"},
		Size:  2,
	}

	counter, _, err := prompt.Run()
	utils.Red.Println(errors.New("Prompt err:" + err.Error()))

	if counter == 0 {
		return "cluster"
	} else if counter == 1 {
		return "namespace"
	}

	return utils.DefaultMode
}

// GetAgentDetails take details of agent as input
func GetAgentDetails(mode string, pid string, c types.Credentials, kubeconfig *string) (types.Agent, error) {
	var (
		newAgent types.Agent
		err      error
	)
	// Get agent name as input
	utils.White_B.Println("\nEnter the details of the agent")
	// Label for goto statement in case of invalid agent name

AGENT_NAME:
	prompt := promptui.Prompt{
		Label: "What's the Agent Name?",
	}

	newAgent.AgentName, err = prompt.Run()
	utils.Red.Println(errors.New("Prompt err:" + err.Error()))

	if newAgent.AgentName == "" {
		utils.Red.Println("⛔ Agent name cannot be empty. Please enter a valid name.")
		goto AGENT_NAME
	}

	// Check if agent with the given name already exists
	agent, err := apis.GetAgentList(c, pid)
	if err != nil {
		return types.Agent{}, err
	}

	var isAgentExist = false
	for i := range agent.Data.GetAgent {
		if newAgent.AgentName == agent.Data.GetAgent[i].AgentName {
			utils.White_B.Println(agent.Data.GetAgent[i].AgentName)
			isAgentExist = true
		}
	}

	if isAgentExist {
		goto AGENT_NAME
	}

	// Get agent description as input
	prompt = promptui.Prompt{
		Label: "Add your agent description",
	}

	newAgent.Description, err = prompt.Run()
	utils.Red.Println(errors.New("Prompt err:" + err.Error()))

	nodeSelector := promptui.Select{
		Label: "Do you want NodeSelectors added to the agent deployments?",
		Items: []string{"Yes", "No"},
	}

	counter, _, err := nodeSelector.Run()
	utils.Red.Println(errors.New("Prompt err:" + err.Error()))

	if counter == 0 {
		prompt := promptui.Prompt{
			Label: "Add your NodeSelector(s) (Format: key1=value1,key2=value2)",
		}
		newAgent.NodeSelector, err = prompt.Run()
		if ok := utils.CheckKeyValueFormat(newAgent.NodeSelector); !ok {
			os.Exit(1)
		}
		utils.Red.Println(errors.New("Prompt err:" + err.Error()))
	}

	toleration := promptui.Select{
		Label: "Do you want Tolerations added in the agent deployments??",
		Items: []string{"Yes", "No"},
	}

	counter, _, err = toleration.Run()
	utils.Red.Println(errors.New("Prompt err:" + err.Error()))

	if counter == 0 {
		prompt := promptui.Prompt{
			Label: "Add the toleration count",
		}

		result, err := prompt.Run()
		utils.Red.Println(errors.New("Prompt err:" + err.Error()))

		nts, err := strconv.Atoi(result)
		utils.PrintError(err)

		str := "["
		for tol := 0; tol < nts; tol++ {
			str += "{"

			utils.White_B.Print("\nToleration count: ", tol+1)

			prompt.Label = "TolerationSeconds: (Press Enter to ignore)"
			ts, err := prompt.Run()
			utils.Red.Println(errors.New("Prompt err:" + err.Error()))

			prompt.Label = "Operator"
			operator, err := prompt.Run()
			utils.Red.Println(errors.New("Prompt err:" + err.Error()))

			if operator != "" {
				str += "operator : \\\"" + operator + "\\\" "
			}

			prompt.Label = "Effect"
			effect, err := prompt.Run()
			utils.Red.Println(errors.New("Prompt err:" + err.Error()))

			if effect != "" {
				str += "effect: \\\"" + effect + "\\\" "
			}

			if ts != "" {
				str += "tolerationSeconds: " + ts + " "
			}

			prompt.Label = "Key"
			key, err := prompt.Run()
			utils.Red.Println(errors.New("Prompt err:" + err.Error()))

			if key != "" {
				str += "key: \\\"" + key + "\\\" "
			}

			prompt.Label = "Value"
			value, err := prompt.Run()
			utils.Red.Println(errors.New("Prompt err:" + err.Error()))

			if key != "" {
				str += "value: \\\"" + value + "\\\" "
			}

			str += " }"
		}
		str += "]"

		newAgent.Tolerations = str

	}

	// Get platform name as input
	newAgent.PlatformName = GetPlatformName(kubeconfig)
	// Set agent type
	newAgent.ClusterType = utils.AgentType
	// Set project id
	newAgent.ProjectId = pid
	// Get namespace
	newAgent.Namespace, newAgent.NsExists = k8s.ValidNs(mode, utils.ChaosAgentLabel, kubeconfig)

	return newAgent, nil
}

func ValidateSAPermissions(mode string, kubeconfig *string) {
	var (
		pems      [2]bool
		err       error
		resources [2]string
	)

	if mode == "cluster" {
		resources = [2]string{"clusterrole", "clusterrolebinding"}
	} else {
		resources = [2]string{"role", "rolebinding"}
	}

	for i, resource := range resources {
		pems[i], err = k8s.CheckSAPermissions(k8s.CheckSAPermissionsParams{Verb: "create", Resource: resource, Print: true}, kubeconfig)
		if err != nil {
			utils.Red.Println(err)
		}
	}

	for _, pem := range pems {
		if !pem {
			utils.Red.Println("\n🚫 You don't have sufficient permissions.\n🙄 Please use a service account with sufficient permissions.")
			os.Exit(1)
		}
	}

	utils.White_B.Println("\n🌟 Sufficient permissions. Installing the Agent...")
}

// Summary display the agent details based on input
func Summary(agent types.Agent, kubeconfig *string) {
	utils.White_B.Printf("\n📌 Summary \nAgent Name: %s\nAgent Description: %s\nPlatform Name: %s\n", agent.AgentName, agent.Description, agent.PlatformName)
	if ok, _ := k8s.NsExists(agent.Namespace, kubeconfig); ok {
		utils.White_B.Println("Namespace: ", agent.Namespace)
	} else {
		utils.White_B.Println("Namespace: ", agent.Namespace, "(new)")
	}

	if k8s.SAExists(k8s.SAExistsParams{Namespace: agent.Namespace, Serviceaccount: agent.ServiceAccount}, kubeconfig) {
		utils.White_B.Println("Service Account: ", agent.ServiceAccount)
	} else {
		utils.White_B.Println("Service Account: ", agent.ServiceAccount, "(new)")
	}

	utils.White_B.Printf("\nInstallation Mode: %s\n", agent.Mode)
}

func ConfirmInstallation() {

	prompt := promptui.Select{
		Label: "Do you want to continue with the above details?",
		Items: []string{"Yes", "No"},
	}

	decision, _, err := prompt.Run()
	utils.Red.Println(errors.New("Prompt err:" + err.Error()))

	if decision == 0 {
		utils.White_B.Println("👍 Continuing agent connection!!")
	} else {
		utils.Red.Println("✋ Exiting agent connection!!")
		os.Exit(1)
	}
}

func CreateRandomProject(userID string, cred types.Credentials) string {
	rand, err := utils.GenerateRandomString(10)
	utils.PrintError(err)

	projectName := cred.Username + "-" + rand

	project, err := apis.CreateProjectRequest(userID, projectName, cred)
	utils.PrintError(err)

	return project.Data.ID
}
