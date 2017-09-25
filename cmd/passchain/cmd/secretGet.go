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
)

// secretGetCmd represents the secretGet command
var secretGetCmd = &cobra.Command{
	Use:   "get",
	Short: "get the value of a secret",
	Long:  `Get a secret, decrypt and print the value.`,
	Run: func(cmd *cobra.Command, args []string) {
		sid := viper.GetString("sid")
		if len(args) > 0 {
			sid = args[0]
		}
		cli := getCli()
		secret, err := cli.GetSecret(sid)
		if err != nil {
			log.Fatal(err)
		}
		encryptedAESKey, ok := secret.Shares[viper.GetString("id")]
		if !ok {
			log.Fatal("no share for us on this secret")
		}
		k := getKey()
		aesKey, err := k.DecryptString(encryptedAESKey)
		if err != nil {
			log.Fatal(err)
		}
		err = secret.Decrypt(aesKey)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(secret.Value)
	},
}

func init() {
	secretCmd.AddCommand(secretGetCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// secretGetCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// secretGetCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
