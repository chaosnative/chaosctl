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
package config

import (
	"github.com/spf13/cobra"
)

// configCmd represents the config command
var ConfigCmd = &cobra.Command{
	Use: "config",
	Short: `It manages multiple litmusportal accounts within a system. 
		Examples(s)
		#set a new account
		chaosctl config set-account  --endpoint "" --access_id "" --access_key ""

		#use an existing account from the config file
		chaosctl config use-account  --endpoint "" --access_id ""

		#get all accounts in the config file
		chaosctl config get-accounts
		
		#view the config file
		chaosctl config view
		
		Note: The default location of the config file is $HOME/.chaosconfig, and can be overridden by a --config flag
		`,
}
