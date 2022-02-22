package hubgen

import "github.com/spf13/cobra"

var HubgenCmd = &cobra.Command{
	Use: "hubgen",
	Short: `Create customized ChaosHub for CLC and CLE.
		Examples:
		#Generate a ChaosHub
		chaosctl hubgen generate --hubname="my-hub-name" --import-path="./my-directory" --export-path="./my-export-path" 

		#Clone a ChaosHub using repo url
		chaosctl hubgen get-charts --hub-path="./my-directory" --repo-url="my-repo-url"
	`,
}
