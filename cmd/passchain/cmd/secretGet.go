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
