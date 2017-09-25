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
)

// secretDelCmd represents the secretDel command
var secretDelCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete a secret",
	Long:  `Delete a secret.`,
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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// secretDelCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// secretDelCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
