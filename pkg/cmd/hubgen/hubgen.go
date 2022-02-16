package hubgen

import "github.com/spf13/cobra"

var HubgenCmd = &cobra.Command{
	Use: "hubgen",
	Short: `Create customized ChaosHub for CLC and CLE.
		Examples:
		#Get local hub path
		chaosctl hubgen import --hubpath="./chaos-hub" 

		#export a hub
		chaosctl hubgen generate --exportpath="./my-directory"
	`,
}
