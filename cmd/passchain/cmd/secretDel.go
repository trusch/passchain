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
)

// secretDelCmd represents the secretDel command
var secretDelCmd = &cobra.Command{
	Use:     "delete",
	Aliases: []string{"remove", "del"},
	Short:   "delete a secret",
	Long:    `Delete a secret.`,
	Run: func(cmd *cobra.Command, args []string) {
		cli := getCli()
		id := viper.GetString("sid")
		if len(args) > 0 {
			id = args[0]
		}
		if id == "" {
			log.Fatal("you must specify --sid")
		}
		if err := cli.DelSecret(id); err != nil {
			log.Fatal(err)
		}
		log.Print("successfully deleted secret ", id)
	},
}

func init() {
	secretCmd.AddCommand(secretDelCmd)
}
