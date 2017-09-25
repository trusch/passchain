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
	"github.com/trusch/passchain/crypto"
	"github.com/trusch/passchain/state"
)

// createAccountCmd represents the createAccount command
var createAccountCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"add"},
	Short:   "create a account",
	Long:    `Here you can create a new account`,
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
		fmt.Printf("export PASSCHAIN_ID=%v\n", id)
		fmt.Printf("export PASSCHAIN_PUBLIC_KEY=%v\n", cli.Key.GetPubString())
		fmt.Printf("export PASSCHAIN_PRIVATE_KEY=%v\n", cli.Key.GetPrivString())
	},
}

func init() {
	accountCmd.AddCommand(createAccountCmd)
}
