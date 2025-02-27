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
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/ed-henrique/tatu/cli/internal/endpoints"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	isFile bool

	errNoServerInConfig = errors.New("No server was added to config. Use the command below to add one:\n\ntatu server add URL")
)

// storeCmd represents the store command
func (cli *CLI) NewStoreCmd() *cobra.Command {
	storeCmd := &cobra.Command{
		Use:   "store SECRET",
		Short: "Store your secret",
		Long:  `Store your secret in encrypted format.`,
		Example: `tatu store "secret"
tatu store -f secret.txt
cat secret.txt | tatu store -
tatu store - < secret.txt`,
		RunE: func(cmd *cobra.Command, args []string) error {
			var (
				input  []byte
				secret []byte
				err    error

				// Only pipes if the arg passed is '-' (Might change this to pipe when len(args) == 0 in the
				// future)
				isPiped = args[0] == "-"
			)

			// Reference https://github.com/spf13/cobra/issues/1749 for more info on this strategy to gather
			// pipe input. There might be changes regarding it in the future.
			inputReader := cmd.InOrStdin()

			if isPiped {
				input, err = io.ReadAll(inputReader)
				if err != nil {
					return err
				}
			} else {
				input = []byte(args[0])
			}

			if isFile {
				var secretFilepath string
				secretFilepath = strings.TrimSpace(string(input))
				absFilepath, err := filepath.Abs(secretFilepath)
				if err != nil {
					return err
				}

				secret, err = os.ReadFile(absFilepath)
				if err != nil {
					return err
				}
			} else {
				secret = input
			}

			// Encode secret as base64 before sending it
			base64Secret := base64.URLEncoding.EncodeToString(secret)
			body, err := json.Marshal(struct {
				Secret string `json:"secret,omitempty"`
			}{
				Secret: base64Secret,
			})

			// Get "server" from config. This config name will change in the future, since multi-server
			// configs are expected
			sr := strings.NewReader(string(body))
			serverURL := viper.GetString("server")
			if serverURL == "" {
				if flaggedServer == "" {
					cmd.SilenceUsage = true // No help message after this error message
					return errNoServerInConfig
				}

				serverURL = flaggedServer
			}

			r, err := http.Post(endpoints.Join(serverURL, endpoints.Secrets), "application/json", sr)
			if err != nil {
				return err
			}

			if r.StatusCode != http.StatusCreated {
				rBody, err := io.ReadAll(r.Body)
				if err != nil {
					return err
				}
				defer r.Body.Close()

				fmt.Fprintln(cmd.ErrOrStderr(), string(rBody))
				return nil
			}

			fmt.Fprintln(cmd.OutOrStdout(), "Secret was added.")
			return nil
		},
		Args: cobra.ExactArgs(1),
	}

	storeCmd.Flags().BoolVarP(&isFile, "file", "f", false, "SECRET is a file path")

	return storeCmd
}
