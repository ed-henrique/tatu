package cmd

import (
	"bytes"
	"io"
	"testing"
)

func TestStoreCmd(t *testing.T) {
	b := new(bytes.Buffer)
	expected := "abc\n"

	rootCmd.SetOut(b)
	rootCmd.SetArgs([]string{"store", "abc"})
	err := rootCmd.Execute()
	if err != nil {
		t.Fatal(err)
	}

	out, err := io.ReadAll(b)
	if err != nil {
		t.Fatal(err)
	}

	if string(out) != expected {
		t.Errorf("expected %s got %s", expected, string(out))
	}
}
