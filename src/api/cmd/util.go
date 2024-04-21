package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"prospector/middleware"
	"strings"
)

func CheckProspectorReachability(address string) bool {
	healthURL := address + "/api/health"
	res, err := http.Get(healthURL)
	if err != nil {
		return false
	}

	if res.StatusCode != 200 {
		return false
	}

	return true
}

func ProspectorAuth(address, username, password string) (string, error) {
	res, err := CmdPost(address+"/api/login", fmt.Sprintf(`{"username":"%s","password":"%s"}`, username, password))
	if err != nil {
		return "", err
	}

	if res.StatusCode != 200 {
		return "", fmt.Errorf("the server returned an error. please try again")
	}

	var authResponse middleware.AuthSucess
	err = json.NewDecoder(res.Body).Decode(&authResponse)
	if err != nil {
		return "", err
	}

	return authResponse.Token, nil
}

func CmdGet(url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("accept", "application/json")

	token := getStoredToken()

	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func CmdDelete(url string, body string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodDelete, url, strings.NewReader(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("accept", "application/json")

	token := getStoredToken()

	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func CmdPost(url string, body string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("accept", "application/json")

	token := getStoredToken()

	if !strings.Contains(url, "/api/login") {
		if token == "" {
			return nil, fmt.Errorf("no token found")
		}
		req.Header.Set("Authorization", "Bearer "+token)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func CmdPut(url string, body string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPut, url, strings.NewReader(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("accept", "application/json")

	token := getStoredToken()

	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func getStoredToken() string {
	b, err := os.ReadFile(os.Getenv("HOME") + "/.prospector_token")
	if err != nil {
		return ""
	}

	return string(b)
}
