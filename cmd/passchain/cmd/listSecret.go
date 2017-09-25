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
)

// listSecretCmd represents the listSecret command
var listSecretCmd = &cobra.Command{
	Use:   "list",
	Short: "list secrets",
	Long:  `List secrets and decrypt if possible.`,
	Run: func(cmd *cobra.Command, args []string) {
		cli := getCli()
		secrets, err := cli.ListSecrets()
		if err != nil {
			log.Fatal(err)
		}
		for _, s := range secrets {
			if encryptedKey, ok := s.Shares[cli.AccountID]; ok {
				key, e := cli.Key.DecryptString(encryptedKey)
				if e != nil {
					log.Fatal(e)
				}
				if e = s.Decrypt(key); e != nil {
					log.Fatal(e)
				}
			}
		}
		print(secrets)
	},
}

func init() {
	secretCmd.AddCommand(listSecretCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listSecretCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listSecretCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
