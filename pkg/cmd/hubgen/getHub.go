package hubgen

import (
	"errors"
	"os"

	"github.com/chaosnative/chaosctl/pkg/utils"
	"github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"
)

var getHubCmd = &cobra.Command{
	Use:   "get-charts",
	Short: "Imports the local ChaosHub with it's path and generates the configured ChaosHub",
	Long:  "Imports the local ChaosHub with it's path and generates the configured ChaosHub",
	Run: func(cmd *cobra.Command, args []string) {
		hubPath, err := cmd.Flags().GetString("hub-path")
		utils.PrintError(err)

		if len(hubPath) == 0 {
			utils.PrintError(errors.New("path to clone ChaosHub not provided"))
		}
		err = validateDir(hubPath)
		if err != nil {
			utils.PrintError(err)
		}

		repoUrl, err := cmd.Flags().GetString("repo-url")
		utils.PrintError(err)
		if len(repoUrl) == 0 {
			repoUrl = "https://github.com/litmuschaos/chaos-charts"
		}

		err = cloneRepo(repoUrl, hubPath)
		if err != nil {
			utils.PrintError(err)
		}
		utils.White_B.Print("\n🎉 ChaosHub cloned successfully at " + hubPath)

	},
}

//validateDir is used to validate the directory using the directory path
func validateDir(path string) error {
	_, err := os.Stat(path)
	if err != nil {
		return err
	}
	if os.IsNotExist(err) {
		return err
	}
	return nil
}

//cloneRepo is used to clone a git repo using repo URL and directory path
func cloneRepo(repoUrl string, clonePath string) error {
	_, err := git.PlainClone(clonePath, false, &git.CloneOptions{
		URL:      repoUrl,
		Progress: os.Stdout,
	})
	if err != nil {
		return err
	}
	return nil
}

func init() {
	HubgenCmd.AddCommand(getHubCmd)
	getHubCmd.Flags().String("hub-path", "", "Path to clone the default ChaosHub")
	getHubCmd.Flags().String("repo-url", "", "Repository URL of ChaosHub to clone")
}
