/*
 * Copyright (C) 2017 Tino Rusch
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of
 * this software and associated documentation files (the "Software"), to deal in
 * the Software without restriction, including without limitation the rights to
 * use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
 * the Software, and to permit persons to whom the Software is furnished to do so,
 * subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
 * FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
 * COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
 * IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
 * CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 */

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
		encryptedAESKey, ok := secret.Shares[cli.AccountID]
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
		if viper.GetBool("owner") {
			secret.Owners[with] = true
		}
		err = cli.UpdateSecret(secret)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	secretCmd.AddCommand(secretShareCmd)
	secretShareCmd.Flags().String("with", "", "who to share with")
	secretShareCmd.Flags().Bool("owner", false, "share owner rights (read only if false)")
	viper.BindPFlags(secretShareCmd.Flags())
}
