# ChaosCTL
The CLC command-line tool, chaosctl, allows you to manage CLC's chaos delegate plane. You can use chaosctl to create chaos delegates, project, and manage multiple CLC accounts. 

## Usage
* For more information including a complete list of chaosctl operations, see the chaosctl reference documentation. 
* Non-Interactive mode: <a href="https://github.com/chaosnative/chaosctl/blob/master/Usage.md">Click here</a>
* Interactive mode: <a href="https://github.com/chaosnative/chaosctl/blob/master/Usage_interactive.md">Click here</a>

## Requirements

The chaosctl CLI requires the following things:

- kubeconfig - chaosctl needs the kubeconfig of the k8s cluster where we need to connect CLC chaos delegates. The CLI currently uses the default path of kubeconfig i.e. `~/.kube/config`.
- kubectl- chaosctl is using kubectl under the hood to apply the manifest. To install kubectl, follow:  [kubectl](https://kubernetes.io/docs/tasks/tools/#kubectl)


## Installation

To install the latest version of chaosctl follow the below steps:

<table>
  <th>Platforms</th>
  <th>main(Unreleased)</th>
  <th>v0.1.0</th>
  <tr>
    <td>chaosctl-darwin-amd64 (MacOS)</td>
    <td><a href="https://chaosctl.chaosnative.com/chaosctl-darwin-amd64-main.tar.gz">Click here</a></td>
    <td><a href="https://chaosctl.chaosnative.com/chaosctl-darwin-amd64-v0.1.0.tar.gz">Click here</a></td>
  </tr>
  <tr>
    <td>chaosctl-linux-386</td>
    <td><a href="https://chaosctl.chaosnative.com/chaosctl-linux-386-main.tar.gz">Click here</a></td>
    <td><a href="https://chaosctl.chaosnative.com/chaosctl-linux-386-v0.1.0.tar.gz">Click here</a></td>
  </tr>
  <tr>
    <td>chaosctl-linux-amd64</td>
    <td><a href="https://chaosctl.chaosnative.com/chaosctl-linux-amd64-main.tar.gz">Click here</a></td>
    <td><a href="https://chaosctl.chaosnative.com/chaosctl-linux-amd64-v0.1.0.tar.gz">Click here</a></td>
  </tr>
  <tr>
    <td>chaosctl-linux-arm</td>
    <td><a href="https://chaosctl.chaosnative.com/chaosctl-linux-arm-main.tar.gz">Click here</a></td>
    <td><a href="https://chaosctl.chaosnative.com/chaosctl-linux-arm-v0.1.0.tar.gz">Click here</a></td>
  </tr>
  <tr>
    <td>chaosctl-linux-arm64</td>
    <td><a href="https://chaosctl.chaosnative.com/chaosctl-linux-arm64-main.tar.gz">Click here</a></td>
    <td><a href="https://chaosctl.chaosnative.com/chaosctl-linux-arm64-v0.1.0.tar.gz">Click here</a></td>
  </tr>
  <tr>
    <td>chaosctl-windows-386</td>
    <td><a href="https://chaosctl.chaosnative.com/chaosctl-windows-386-main.tar.gz">Click here</a></td>
    <td><a href="https://chaosctl.chaosnative.com/chaosctl-windows-386-v0.1.0.tar.gz">Click here</a></td>
  </tr>
   <tr>
    <td>chaosctl-windows-amd64</td>
    <td><a href="https://chaosctl.chaosnative.com/chaosctl-windows-amd64-main.tar.gz">Click here</a></td>
    <td><a href="https://chaosctl.chaosnative.com/chaosctl-windows-amd64-v0.1.0.tar.gz">Click here</a></td>
  </tr>
  <tr>
    <td>chaosctl-windows-arm</td>
    <td><a href="https://chaosctl.chaosnative.com/chaosctl-windows-arm-main.tar.gz">Click here</a></td>
    <td><a href="https://chaosctl.chaosnative.com/chaosctl-windows-arm-v0.1.0.tar.gz">Click here</a></td>
  </tr>
</table>

### Linux/MacOS

* Extract the binary

```shell
tar -zxvf chaosctl-<OS>-<ARCH>-<VERSION>.tar.gz
```

* Provide necessary permissions

```shell
chmod +x chaosctl
```

* Move the chaosctl binary to /usr/local/bin/chaosctl. Note: Make sure to use root user or use sudo as a prefix

```shell
mv chaosctl /usr/local/bin/chaosctl
```

* You can run the chaosctl command in Linux/macOS:

```shell
chaosctl <command> <subcommand> <subcommand> [options and parameters]
```

### Windows

* Extract the binary from the zip using WinZip or any other extraction tool.

* You can run the chaosctl command in windows:

```shell
chaosctl.exe <command> <subcommand> <subcommand> [options and parameters]
```

* To check the version of the chaosctl:

```shell
chaosctl version
```

----
