# Usage: Chaosctl (Interactive mode)

### chaosctl Syntax

`chaosctl` has a syntax to use as follows:

```shell
chaosctl [command] [TYPE] [flags]
```

- Command: refers to what you do want to perform (create, get and config)
- Type: refers to the feature type you are performing a command against (agent, project etc.)
- Flags: It takes some additional information for resource operations. For example, `--installation-mode` allows you to specify an installation mode.

chaosctl is using the `.chaosconfig` config file to manage multiple accounts

1. If the --config flag is set, then only the given file is loaded. The flag may only be set once and no merging takes place.
2. Otherwise, the ${HOME}/.chaosconfig file is used, and no merging takes place.

chaosctl supports both interactive and non-interactive(flag based) modes.

> Only `chaosctl create agent` command needs --non-interactive flag, other commands don't need this flag to be in non-interactive mode. If mandatory flags aren't passed, then chaosctl takes input in an interactive mode.

### Steps to create an agent

- To setup an account with chaosctl

```shell
chaosctl config set-account
```

Next, you need to enter CLC/CLE details to login into your account. Fields to be filled in:

**What's the product name:** Select the product name .

> Example, https://preview.litmuschaos.io/

**AccessID:** What's the AccessID?: <br />
**AccessKey:** What's the AccessKey?:

```
? What's the product name?: 
  ▸ ChaosNative Cloud
    ChaosNative Enterprise
    
What's the AccessID?: Raj60163RjxQE
What's the AccessKey?: ***************

account.accessID/admin configured
```

- To create an agent in a cluster mode

```shell
chaosctl create agent
```

There will be a list of existing projects displayed on the terminal. Select the desired project by entering the sequence number indicated against it.

```
? Select a project from the list: 
  ▸ Raj60163's project
```

Next, select the installation mode based on your requirement by entering the sequence number indicated against it.

It can install an agent in two different modes.

- cluster mode: With this mode, the agent can run the chaos in any namespace. It installs appropriate cluster roles and cluster role bindings to achieve this mode.

- namespace mode: With this mode, the agent can run the chaos in its namespace. It installs appropriate roles and role bindings to achieve this mode.

Note: With namespace mode, the user needs to create the namespace to install the agent as a prerequisite.

```
? What's the installation mode?: 
  ▸ Cluster
    Namespace

🏃 Running prerequisites check....
🔑 clusterrole ✅
🔑 clusterrolebinding ✅
🌟 Sufficient permissions. Installing the Agent...

```

Next, enter the details of the new agent.

Fields to be filled in <br />

<table>
    <th>Field</th>
    <th>Description</th>
    <tr>
        <td>Agent Name:</td>
        <td>Enter a name of the agent which needs to be unique across the project</td>
    </tr>
    <tr>
        <td>Agent Description:</td>
        <td>Fill in details about the agent</td>
    </tr>
    <tr>
        <td>Node Selector:</td>
        <td>To deploy the agent on a particular node based on the node selector labels</td>
    </tr>
    <tr>
        <td>Platform Name:</td>
        <td>Enter the platform name on which this agent is hosted. For example, AWS, GCP, Rancher etc.</td>
    </tr>
    <tr>
        <td>Enter the namespace:</td>
        <td>You can either enter an existing namespace or enter a new namespace. In cases where the namespace does not exist, chaosctl creates it for you</td>
    </tr>
    <tr>
        <td>Enter service account:</td>
        <td>You can either enter an existing or new service account</td>
    </tr>
</table>

```
Enter the details of the agent
✔ What's the Agent Name?: new-agent

✔ Add your agent description: new-agent█

? Do you want NodeSelectors added to the agent deployments?: 
    Yes
  ▸ No

? Do you want Tolerations added in the agent deployments??: 
    Yes
  ▸ No

? What's your Kubernetes Platform?: 
  ▸ Others
    AWS Elastic Kubernetes Service
    Google Kubernetes Service
    OpenShift
    Rancher

✔ Enter a new or existing namespace [Default: litmus ]: new-ns

✔ Enter a service account [Default: litmus ]:
```

Once, all these steps are implemented you will be able to see a summary of all the entered fields.
After verification of these details, you can proceed with the connection of the agent by entering Y. The process of connection might take up to a few seconds.

```
Enter service account [Default: litmus]:

📌 Summary
Agent Name: New-Agent
Agent Description: This is a new agent
Platform Name: Others
Namespace:  litmus
Service Account:  litmus (new)
Installation Mode: cluster

? Do you want to continue with the above details?: 
  ▸ Yes
    No

👍 Continuing agent connection!!
Applying YAML:
https://preview.litmuschaos.io/api/file/eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjbHVzdGVyX2lkIjoiMDUyZmFlN2UtZGM0MS00YmU4LWJiYTgtMmM4ZTYyNDFkN2I0In0.i31QQDG92X5nD6P_-7TfeAAarZqLvUTFfnAghJYXPiM.yaml

💡 Connecting agent to Litmus Portal.
🏃 Agents are running!!

🚀 Agent Connection Successful!! 🎉
👉 Litmus agents can be accessed here: https://preview.litmuschaos.io/targets
```

#### Verify the new Agent Connection\*\*

To verify, if the connection process was successful you can view the list of connected agents from the Targets section on your ChaosCenter and ensure that the connected agent is in Active State.

---

### Additional commands

- To view the current configuration of `.chaosconfig`, type:

```shell
chaosctl config view
```

**Output:**

```
accounts:
- users:
  - expires_in: "1626897027"
    token: eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MjY4OTcwMjcsInJvbGUiOiJhZG1pbiIsInVpZCI6ImVlODZkYTljLTNmODAtNGRmMy04YzQyLTExNzlhODIzOTVhOSIsInVzZXJuYW1lIjoiYWRtaW4ifQ.O_hFcIhxP4rhyUN9NEVlQmWesoWlpgHpPFL58VbJHnhvJllP5_MNPbrRMKyFvzW3hANgXK2u8437u
    username: admin
  - expires_in: "1626944602"
    token: eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MjY5NDQ2MDIsInJvbGUiOiJ1c2VyIiwidWlkIjoiNjFmMDY4M2YtZWY0OC00MGE1LWIzMjgtZTU2ZDA2NjM1MTE4IiwidXNlcm5hbWUiOiJyYWoifQ.pks7xjkFdJD649RjCBwQuPF1_QMoryDWixSKx4tPAqXI75ns4sc-yGhMdbEvIZ3AJSvDaqTa47XTC6c8R
    username: litmus-user
  endpoint: https://preview.litmuschaos.io
apiVersion: v1
current-account: https://preview.litmuschaos.io
current-user: litmus-user
kind: Config
```

- To get an overview of the accounts available within `.chaosconfig`, use the `config get-accounts` command:

```shell
chaosctl config get-accounts
```

**Output:**

```
CURRENT  ENDPOINT                         ACCESSID  EXPIRESIN
         https://preview.litmuschaos.io   admin     2021-07-22 01:20:27 +0530 IST
*        https://preview.litmuschaos.io   raj       2021-07-22 14:33:22 +0530 IST
```

- To alter the current account use the `use-account` command:

```shell
chaosctl config use-account

? What's the product name?: 
  ▸ ChaosNative Cloud
    ChaosNative Enterprise
    
What's the AccessID?: Raj60163RjxQE
```

- To create a project, apply the following command :

```shell
chaosctl create project

Enter a project name: new
```

- To view all the projects with the user, use the `get projects` command.

```shell
chaosctl get projects
```

**Output:**

```
PROJECT ID                                PROJECT NAME       CREATEDAT
50addd40-8767-448c-a91a-5071543a2d8e      Developer Project  2021-07-21 14:38:51 +0530 IST
7a4a259a-1ae5-4204-ae83-89a8838eaec3      DevOps Project     2021-07-21 14:39:14 +0530 IST
```

- To get an overview of the agents available within a project, issue the following command.

```shell
chaosctl get agents

Enter the Project ID: 50addd40-8767-448c-a91a-5071543a2d8e
```

**Output:**

```
AGENTID                                AGENTNAME          STATUS
55ecc7f2-2754-43aa-8e12-6903e4c6183a   agent-1            ACTIVE
13dsf3d1-5324-54af-4g23-5331g5v2364f   agent-2            INACTIVE
```

For more information related to flags, Use `chaosctl --help`.
