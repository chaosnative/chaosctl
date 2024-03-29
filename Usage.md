### chaosctl Syntax

`chaosctl` has a syntax to use as follows:

```shell
chaosctl [command] [TYPE] [flags]
```

- Command: refers to what you do want to perform (connect, create, get and config)
- Type: refers to the feature type you are performing a command against (chaos-delegate, project etc.)
- Flags: It takes some additional information for resource operations. For example, `--installation-mode` allows you to specify an installation mode.

chaosctl is using the `.chaosconfig` config file to manage multiple accounts

1. If the --config flag is set, then only the given file is loaded. The flag may only be set once and no merging takes place.
2. Otherwise, the ${HOME}/.chaosconfig file is used, and no merging takes place.

chaosctl supports both interactive and non-interactive(flag based) modes.

> Only `chaosctl connect chaos-delegate` command needs --non-interactive flag, other commands don't need this flag to be in non-interactive mode. If mandatory flags aren't passed, then chaosctl takes input in an interactive mode.

### Installation modes

chaosctl can install a chaos delegate in two different modes.

- cluster mode: With this mode, the chaos delegate can run the chaos in any namespace. It installs appropriate cluster roles and cluster role bindings to achieve this mode. It can be enabled by passing a flag `--installation-mode=cluster`

- namespace mode: With this mode, the chaos delegate can run the chaos in its namespace. It installs appropriate roles and role bindings to achieve this mode. It can be enabled by passing a flag `--installation-mode=namespace`

Note: With namespace mode, the user needs to create the namespace to install the chaos delegate and user must have the admin privileges to setup [CRDs](https://github.com/chaosnative/hce-charts/blob/main/k8s-manifests/ci/hce-crds.yaml) as a prerequisite.

#### Prerequisite steps(For namespace mode)

- Create namespace to install the chaos delegate

```shell
kubectl create ns <namespace_name>
```

- Setup CRDs

```shell
kubectl apply -f https://raw.githubusercontent.com/litmuschaos/litmus/master/litmus-portal/litmus-portal-crds.yml
```

### Minimal steps to connect a chaos delegate

- To setup an account with chaosctl

```shell
chaosctl config set-account --endpoint="" --access_id="" --access_key=""
```

- To connect a chaos delegate with an existing project
  > Note: To get `project-id`. Apply `chaosctl get projects`

```shell
chaosctl connect chaos-delegate --name="" --project-id="" --non-interactive
```

### Flags for `connect chaos-delegate` command

<table>
<tr>
    <th>Flag</th>
    <th>Short Flag</th>
    <th>Type</th>
    <th>Description</th>
    <tr>
        <td>--description</td>
        <td></td>
        <td>String</td>
        <td>Set the chaos delegate description (default "---")</td>
    </tr>
    <tr>
        <td>--name</td>
        <td></td>
        <td>String</td>
        <td>Set the chaos-delegate-type to external for external chaos delegates | Supported=external/internal (default "external")</td>
    </tr>
        <tr>
        <td>--skip-ssl</td>
        <td></td>
        <td>Boolean</td>
        <td>Set whether chaos delegate will skip ssl/tls check (can be used for self-signed certs, if cert is not provided in portal) (default false)</td>
    </tr>
    <tr>
        <td>--chaos-delegate-type</td>
        <td></td>
        <td>String</td>
        <td>Set the chaos-delegate-type to external for external chaos delegates | Supported=external/internal (default "external")</td>
    </tr>
    <tr>
        <td>--installation-mode</td>
        <td></td>
        <td>String</td>
        <td>Set the installation mode for the kind of chaos delegate | Supported=cluster/namespace (default "cluster")</td>
    </tr>
    <tr>
        <td>--kubeconfig</td>
        <td>-k</td>
        <td>String</td>
        <td>Set to pass kubeconfig file if it is not in the default location ($HOME/.kube/config)</td>
    </tr>
    <tr>
        <td>--namespace</td>
        <td></td>
        <td>String</td>
        <td>Set the namespace for the chaos delegate installation (default "litmus")</td>
    </tr>
    <tr>
        <td>--node-selector</td>
        <td></td>
        <td>String</td>
        <td>Set the node-selector for chaos delegate components | Format: key1=value1,key2=value2)
    </tr>
    <tr>
        <td>--non-interactive</td>
        <td>-n</td>
        <td>String</td>
        <td>Set it to true for non interactive mode | Note: Always set the boolean flag as --non-interactive=Boolean</td>
    </tr>
    <tr>
        <td>--ns-exists</td>
        <td></td>
        <td>Boolean</td>
        <td>Set the --ns-exists=false if the namespace mentioned in the --namespace flag is not existed else set it to --ns-exists=true | Note: Always set the boolean flag as --ns-exists=Boolean</td>
    </tr>
    <tr>
        <td>--platform-name</td>
        <td></td>
        <td>String</td>
        <td>Set the platform name. Supported- AWS/GKE/Openshift/Rancher/Others (default "Others")</td>
    </tr>
    <tr>
        <td>--sa-exists</td>
        <td></td>
        <td>Boolean</td>
        <td>Set the --sa-exists=false if the service-account mentioned in the --service-account flag is not existed else set it to --sa-exists=true | Note: Always set the boolean flag as --sa-exists=Boolean"</td>
    </tr>
    <tr>
        <td>--service-account</td>
        <td></td>
        <td>String</td>
        <td>Set the service account to be used by the chaos delegate (default "litmus")</td>
    </tr>
    <tr>
        <td>--config</td>
        <td></td>
        <td>String</td>
        <td>config file (default is $HOME/.chaosconfig)</td>
    </tr>
</table>

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

- To alter the current account use the `use-account` command with the --endpoint and --access_id flags:

```shell
chaosctl config use-account --endpoint="" --access_id=""
```

- To create a project, apply the following command with the `--name` flag:

```shell
chaosctl create project --name=""
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

- To get an overview of the chaos delegates available within a project, issue the following command.

```shell
chaosctl get chaos-delegates --project-id=""
```

**Output:**

```
CHAOS DELEGATE ID                                CHAOS DELEGATE NAME          STATUS
55ecc7f2-2754-43aa-8e12-6903e4c6183a   chaos-delegate-1            ACTIVE
13dsf3d1-5324-54af-4g23-5331g5v2364f   chaos-delegate-2            INACTIVE
```

For more information related to flags, Use `chaosctl --help`.

---
