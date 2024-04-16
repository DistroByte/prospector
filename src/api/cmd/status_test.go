package cmd

import (
	"bytes"
	"testing"
)

func TestExecuteStatusCommand(t *testing.T) {

	actual := new(bytes.Buffer)

	statusCmd.SetOut(actual)
	statusCmd.SetErr(actual)

	httpServer := createTestServer(200, "")

	statusCmd.Flags().Set("address", httpServer.URL)
	defer httpServer.Close()
	statusCmd.Run(statusCmd, []string{"my-job"})

	expected := ``

	if actual.String() != expected {
		t.Errorf("Expected %s but got %s", expected, actual.String())
	}
}

func TestExecuteStatusCommandError(t *testing.T) {

	actual := new(bytes.Buffer)

	statusCmd.SetOut(actual)
	statusCmd.SetErr(actual)
	statusCmd.Run(statusCmd, []string{"my-job"})

	expected := ``

	if actual.String() != expected {
		t.Errorf("Expected %s but got %s", expected, actual.String())
	}
}

func TestExecuteStatusCommandDownstreamError(t *testing.T) {

	actual := new(bytes.Buffer)

	statusCmd.SetOut(actual)
	statusCmd.SetErr(actual)

	httpServer := createTestServer(500, "")

	statusCmd.Flags().Set("address", httpServer.URL)
	defer httpServer.Close()
	statusCmd.Run(statusCmd, []string{"my-job"})

	expected := ``

	if actual.String() != expected {
		t.Errorf("Expected %s but got %s", expected, actual.String())
	}
}
