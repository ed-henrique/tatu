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
package cmds

import (
	"bytes"
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ed-henrique/tatu/internal/server"
	"github.com/spf13/viper"
)

func TestStoreCmd(t *testing.T) {
	testtable := []struct {
		name     string
		args     []string
		expected string
		reader   io.Reader
	}{
		{
			name:     "store string abc",
			args:     []string{"store", "abc"},
			expected: "Secret was added.\n",
		},
		{
			name:     "store string abcde from secret.txt",
			args:     []string{"store", "-f", "../../../txt_tests/secret.txt"},
			expected: "Secret was added.\n",
		},
		{
			name:     "store string abcde from pipe",
			args:     []string{"store", "-"},
			expected: "Secret was added.\n",
			reader:   strings.NewReader("abcde"),
		},
	}

	for _, tt := range testtable {
		b := new(bytes.Buffer)
		s := server.New()
		s.Routes()

		ts := httptest.NewServer(s.Mux)
		defer ts.Close()

		v := viper.New()
		v.Set("server", ts.URL)

		cli := New(
			WithConfig(v),
			WithHTTPClient(ts.Client()),
		)

		cli.root.SetIn(tt.reader)
		cli.root.SetOut(b)
		cli.root.SetArgs(tt.args)

		cli.Execute()

		out, err := io.ReadAll(b)
		if err != nil {
			t.Fatal(err)
		}

		if string(out) != tt.expected {
			t.Errorf("expected %s got %s", tt.expected, string(out))
		}
	}
}
