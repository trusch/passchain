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
	"github.com/trusch/passchain/state"
)

var secretData string

// secretAddCmd represents the secretAdd command
var secretAddCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"add"},
	Short:   "create a secret",
	Long:    `Create a secret.`,
	Run: func(cmd *cobra.Command, args []string) {
		sid := viper.GetString("sid")
		if len(args) > 0 {
			sid = args[0]
		}
		data := secretData
		if data == "" && len(args) > 1 {
			data = args[1]
		}
		if sid == "" || data == "" {
			log.Fatal("you must specify --sid and --data")
		}
		ownID := viper.GetString("id")
		s := &state.Secret{
			ID:     sid,
			Value:  data,
			Shares: make(map[string]string),
			Owners: map[string]bool{
				ownID: true,
			},
		}
		aesKey, err := s.Encrypt()
		if err != nil {
			log.Fatal(err)
		}
		k := getKey()
		encryptedAESKey, err := k.EncryptToString(aesKey)
		if err != nil {
			log.Fatal(err)
		}
		s.Shares[ownID] = encryptedAESKey
		cli := getCli()
		if err := cli.AddSecret(s); err != nil {
			log.Fatal(err)
		}
		log.Printf("created secret %v", sid)
	},
}

func init() {
	secretCmd.AddCommand(secretAddCmd)
	secretAddCmd.Flags().StringVar(&secretData, "data", "", "secret value")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// secretAddCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// secretAddCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
