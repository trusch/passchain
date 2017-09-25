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
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/trusch/passchain/crypto"
)

// secretShareCmd represents the secretShare command
var secretShareCmd = &cobra.Command{
	Use:   "share",
	Short: "share a secret",
	Long:  `Share a secret with another account.`,
	Run: func(cmd *cobra.Command, args []string) {
		cli := getCli()
		sid := viper.GetString("sid")
		if len(args) > 0 {
			sid = args[0]
		}
		with := viper.GetString("with")
		if len(args) > 1 {
			with = args[1]
		}
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
		acc, err := cli.GetAccount(with)
		if err != nil {
			log.Print("can not find account " + with)
			log.Fatal(err)
		}
		otherKey, err := crypto.NewFromStrings(acc.PubKey, "")
		if err != nil {
			log.Fatal(err)
		}
		otherEncrptedAESKey, err := otherKey.EncryptToString(aesKey)
		if err != nil {
			log.Fatal(err)
		}
		secret.Shares[with] = otherEncrptedAESKey
		err = cli.UpdateSecret(secret)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	secretCmd.AddCommand(secretShareCmd)
	secretShareCmd.Flags().String("with", "", "who to share with")
	viper.BindPFlags(secretShareCmd.Flags())
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// secretShareCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// secretShareCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
