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
	"encoding/json"
	"fmt"
	"log"
	"os"

	yaml "gopkg.in/yaml.v2"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/trusch/passchain/client"
	"github.com/trusch/passchain/crypto"
)

var cfgFile string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "passchain",
	Short: "A blockchain based password sharing application",
	Long: `Passchain lets you create and share secrets with your team.

It is based on the tendermint blockchain and handles accounts and secrets.
	`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cli.yaml)")
	RootCmd.PersistentFlags().String("public-key", "", "public key")
	RootCmd.PersistentFlags().String("private-key", "", "private key")
	RootCmd.PersistentFlags().String("id", "", "id of account")
	RootCmd.PersistentFlags().String("format", "yaml", "output format")

	viper.BindPFlags(RootCmd.PersistentFlags())
	viper.BindEnv("id", "PASSCHAIN_ID")
	viper.BindEnv("public-key", "PASSCHAIN_PUBLIC_KEY")
	viper.BindEnv("private-key", "PASSCHAIN_PRIVATE_KEY")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".cli" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".cli")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func getKey() *crypto.Key {
	k, err := crypto.NewFromStrings(viper.GetString("public-key"), viper.GetString("private-key"))
	if err != nil {
		log.Print("generate ephemeral key...")
		k, err = crypto.CreateKeyPair()
		if err != nil {
			log.Fatal(err)
		}
		return k
	}
	return k
}

func getCli() *client.Client {
	return client.NewHTTPClient("http://localhost:46657", getKey(), viper.GetString("id"))
}

func print(data interface{}) {
	format := viper.GetString("format")
	var (
		bs  []byte
		err error
	)
	switch format {
	case "yaml":
		{
			bs, err = yaml.Marshal(data)
		}
	case "json":
		{
			bs, err = json.Marshal(data)
		}
	case "pretty-json":
		{
			bs, err = json.MarshalIndent(data, "", "  ")
		}
	default:
		{
			log.Fatal("wrong format")
		}
	}
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(string(bs))
}
