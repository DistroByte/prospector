package main

import (
	"os"
	"testing"
)

func TestRootNoArgs(t *testing.T) {
	os.Args = []string{"prospector"}
	err := root(os.Args[1:])
	if err == nil {
		t.Error("Expected error")
	}
}

func TestRootUnknownArgs(t *testing.T) {
	os.Args = []string{"prospector", "unknown"}
	err := root(os.Args[1:])
	if err == nil {
		t.Error("Expected error")
	}
}
