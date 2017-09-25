// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/trusch/passchain/crypto"
	"github.com/trusch/passchain/state"
)

// createAccountCmd represents the createAccount command
var createAccountCmd = &cobra.Command{
	Use:   "create",
	Short: "create a account",
	Long:  `Here you can create a new account`,
	Run: func(cmd *cobra.Command, args []string) {
		id := viper.GetString("id")
		if id == "" {
			log.Fatal("you must specify --id")
		}
		if len(args) > 0 {
			id = args[0]
		}
		cli := getCli()
		key, err := crypto.CreateKeyPair()
		if err != nil {
			log.Fatal(err)
		}
		cli.Key = key
		if err := cli.AddAccount(&state.Account{ID: id, PubKey: key.GetPubString()}); err != nil {
			log.Fatal("account creation failed")
		}
		log.Print("successfully created account ", id)
		log.Print("you may wish to set these env variables:")
		fmt.Printf("export PASSCHAIN_ID=%v\n", cli.AccountID)
		fmt.Printf("export PASSCHAIN_PUBLIC_KEY=%v\n", cli.Key.GetPubString())
		fmt.Printf("export PASSCHAIN_PRIVATE_KEY=%v\n", cli.Key.GetPrivString())
	},
}

func init() {
	accountCmd.AddCommand(createAccountCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createAccountCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createAccountCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
