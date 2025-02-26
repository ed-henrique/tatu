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
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// serverAddCmd represents the serverAdd command
var serverAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add new tatu server",
	Long:  `Add new tatu server, saving it on your config.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		viper.Set("server", args[0])
		err := viper.WriteConfig()
		if err != nil {
			return err
		}

		fmt.Fprintln(cmd.OutOrStdout(), "Server was added")
		return nil
	},
	Args: cobra.ExactArgs(1),
}

func init() {
	serverCmd.AddCommand(serverAddCmd)
}
