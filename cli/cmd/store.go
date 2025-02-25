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
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var (
	isFile bool
)

// storeCmd represents the store command
var storeCmd = &cobra.Command{
	Use:   "store SECRET",
	Short: "Store your secret",
	Long:  `Store your secret in encrypted format.`,
	Run: func(cmd *cobra.Command, args []string) {
		if isFile {
			secretFilepath := args[0]
			absFilepath, err := filepath.Abs(secretFilepath)
			if err != nil {
			}

			rawSecret, err := os.ReadFile(absFilepath)
			fmt.Println(rawSecret)
			return
		}

		rawSecret := []byte(args[0])
		fmt.Println(rawSecret)
	},
	Args: cobra.ExactArgs(1),
}

func init() {
	rootCmd.AddCommand(storeCmd)

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	storeCmd.Flags().BoolVarP(&isFile, "file", "f", false, "SECRET is a file path")
}
