package cmd

import (
	"bytes"
	"testing"
)

func TestTemplateCommandExecuteDocker(t *testing.T) {
	actual := new(bytes.Buffer)

	templateCmd.SetOut(actual)
	templateCmd.SetErr(actual)

	httpServer := createTestServer(200, "")
	templateCmd.Flags().Set("address", httpServer.URL)
	templateCmd.Flags().Set("name", "test")
	templateCmd.Flags().Set("components", "2")
	templateCmd.Flags().Set("type", "docker")

	defer httpServer.Close()
	templateCmd.Run(templateCmd, []string{""})

	expected := ``
	if actual.String() != expected {
		t.Errorf("Expected %s but got %s", expected, actual.String())
	}
}

func TestTemplateCommandExecuteVM(t *testing.T) {
	actual := new(bytes.Buffer)

	templateCmd.SetOut(actual)
	templateCmd.SetErr(actual)

	httpServer := createTestServer(200, "")
	templateCmd.Flags().Set("address", httpServer.URL)
	templateCmd.Flags().Set("name", "test")
	templateCmd.Flags().Set("components", "2")
	templateCmd.Flags().Set("type", "vm")

	defer httpServer.Close()
	templateCmd.Run(templateCmd, []string{""})

	expected := ``
	if actual.String() != expected {
		t.Errorf("Expected %s but got %s", expected, actual.String())
	}
}
