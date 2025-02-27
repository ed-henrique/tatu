/*
Copyright © 2025 Eduardo Henrique Freire Machado

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package cmd

import (
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const defaultCfgFile = ".tatu"

var (
	cfgFile       string
	flaggedServer string
)

type CLI struct {
	v    *viper.Viper
	cmds map[string]*cobra.Command
}

func New(v *viper.Viper) *CLI {
	rootCmd := &cobra.Command{
		Use:   "tatu",
		Short: "Secret manager for dev teams",
		Long: `Tatu is a CLI tool designed to safeguard your team's secrets.
Store credentials, API keys, and sensitive data in an encrypted
environment, ensuring security without sacrificing simplicity.`,
		Version: "0.0.0-alpha",
	}

	return &CLI{
		v:    v,
		cmds: map[string]*cobra.Command{"root": rootCmd},
	}
}

func (cli *CLI) AddCommand(parent, name string, cmd func() *cobra.Command) {
	cmdRef := cmd()
	cli.cmds[name] = cmdRef
	cli.cmds[parent].AddCommand(cmdRef)
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func (cli *CLI) Execute() {
	rootCmd := cli.cmds["root"]
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/"+defaultCfgFile+".toml)")
	rootCmd.PersistentFlags().StringVarP(&flaggedServer, "server", "s", "", "host server (defaults to loading from "+defaultCfgFile+".toml)")

	if cfgFile != "" {
		// Use config file from the flag.
		cli.v.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".cli" (without extension).
		cli.v.AddConfigPath(home)
		cli.v.SetConfigType("toml")
		cli.v.SetConfigName(defaultCfgFile)

		// Creates config file if not exists
		err = cli.v.SafeWriteConfig()
		if !strings.Contains(err.Error(), "Already Exists") {
			cobra.CheckErr(err)
		}
	}

	cli.v.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := cli.v.ReadInConfig(); err == nil {
		// fmt.Fprintln(os.Stderr, "Using config file:", cli.v.ConfigFileUsed())
	}

	// Commands
	cli.AddCommand("root", "store", cli.NewStoreCmd)
	cli.AddCommand("root", "server", cli.NewServerCmd)
	cli.AddCommand("server", "server_add", cli.NewServerAddCmd)

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
