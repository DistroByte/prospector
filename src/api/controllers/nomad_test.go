package controllers

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func createTestServer(statusCode int, body string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(statusCode)
		w.Write([]byte(body))
	}))
}

func TestNomadClientGet(t *testing.T) {

	httpServer := createTestServer(200, "")
	defer httpServer.Close()

	TestNomadClient := DefaultNomadClient{
		URL: httpServer.URL,
	}

	_, err := TestNomadClient.Get("/")
	if err != nil {
		t.Errorf("Expected nil but got %s", err)
	}
}

func TestNomadClientGetDownstreamError(t *testing.T) {
	httpServer := createTestServer(500, fmt.Errorf("error").Error())
	defer httpServer.Close()

	TestNomadClient := DefaultNomadClient{
		URL: httpServer.URL,
	}

	_, err := TestNomadClient.Get("/")
	if err == nil {
		t.Errorf("Expected error but got nil")
	}
}

func TestNomadClientPost(t *testing.T) {
	httpServer := createTestServer(200, "")
	defer httpServer.Close()

	TestNomadClient := DefaultNomadClient{
		URL: httpServer.URL,
	}

	body := bytes.NewBufferString(`{"key": "value"}`)
	_, err := TestNomadClient.Post("/", body)
	if err != nil {
		t.Errorf("Expected nil but got %s", err)
	}
}

func TestNomadClientPostError(t *testing.T) {
	httpServer := createTestServer(500, "")
	defer httpServer.Close()

	TestNomadClient := DefaultNomadClient{
		URL: httpServer.URL,
	}

	body := bytes.NewBufferString(`{"key": "value"}`)
	_, err := TestNomadClient.Post("/", body)
	if err == nil {
		t.Errorf("Expected error but got nil")
	}
}

func TestNomadClientDelete(t *testing.T) {
	httpServer := createTestServer(200, "")
	defer httpServer.Close()

	TestNomadClient := DefaultNomadClient{
		URL: httpServer.URL,
	}

	_, err := TestNomadClient.Delete("/")
	if err != nil {
		t.Errorf("Expected nil but got %s", err)
	}
}

func TestNomadClientDeleteError(t *testing.T) {
	httpServer := createTestServer(500, "")
	defer httpServer.Close()

	TestNomadClient := DefaultNomadClient{
		URL: httpServer.URL,
	}

	_, err := TestNomadClient.Delete("/")
	if err == nil {
		t.Errorf("Expected error but got nil")
	}
}

func TestNomadClientDeleteDownstreamError(t *testing.T) {
	httpServer := createTestServer(500, fmt.Errorf("error").Error())
	defer httpServer.Close()

	TestNomadClient := DefaultNomadClient{
		URL: httpServer.URL,
	}

	_, err := TestNomadClient.Delete("/")
	if err == nil {
		t.Errorf("Expected error but got nil")
	}
}
