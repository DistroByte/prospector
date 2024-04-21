package cmd

import (
	"bytes"
	"testing"
)

func TestProjectPreRun(t *testing.T) {
	actual := new(bytes.Buffer)

	projectCmd.SetOut(actual)
	projectCmd.SetErr(actual)

	projectCmd.Flags().Set("address", "https://prospector.ie")
	projectCmd.PersistentPreRun(projectCmd, []string{})

	expected := ``
	if actual.String() != expected {
		t.Errorf("Expected %s but got %s", expected, actual.String())
	}
}
