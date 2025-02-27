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
	"bytes"
	"io"
	"testing"
)

func TestStoreCmd(t *testing.T) {
	// TODO: mock server, preferably even viper
	testtable := []struct {
		name     string
		args     []string
		expected string
		reader   io.Reader
	}{
		{
			name:     "store string abc",
			args:     []string{"--server", "http://localhost:8080", "store", "abc"},
			expected: "Secret was added.\n",
			reader:   nil,
		},
	}

	for _, tt := range testtable {
		b := new(bytes.Buffer)

		rootCmd.SetOut(b)
		rootCmd.SetArgs(tt.args)

		err := rootCmd.Execute()
		if err != nil {
			t.Fatal(err)
		}

		out, err := io.ReadAll(b)
		if err != nil {
			t.Fatal(err)
		}

		if string(out) != tt.expected {
			t.Errorf("expected %s got %s", tt.expected, string(out))
		}
	}
}
