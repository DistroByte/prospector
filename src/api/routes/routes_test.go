package routes

import (
	"prospector/helpers"
	"prospector/middleware"
	"testing"

	"github.com/gin-gonic/gin"
)

func createGinAndController() *gin.Engine {
	r := gin.Default()

	middleware.CreateAuthMiddlewares(r)
	CreateRoutes(r)

	return r
}

func TestCreateRoutes(t *testing.T) {
	createGinAndController()
}

func TestHealth(t *testing.T) {
	r := createGinAndController()

	// Create a request to send to the above route
	req := helpers.PerformRequest(r, "GET", "/api/health")

	// Test that the http status code is 200
	if req.Code != 200 {
		t.Errorf("Expected response code %d. Got %d\n", 200, req.Code)
	}
}

func TestAuthenticateNoUser(t *testing.T) {
	r := createGinAndController()

	req := helpers.PerformAuthRequest(r, "GET", "/api/v1/auth", "", "")

	// Test that the http status code is 401
	if req.Code != 401 {
		t.Errorf("Expected response code %d. Got %d\n", 401, req.Code)
	}
}

// TODO: Re-enable these tests once we have a way to mock the LDAP server

// func TestBasicAuthenticateWithPassword(t *testing.T) {
// 	r := createGinAndController()

// 	req := helpers.PerformAuthRequest(r, "GET", "/api/v1/auth", "admin", "admin")

// 	// Test that the http status code is 200
// 	if req.Code != 200 {
// 		t.Errorf("Expected response code %d. Got %d\n", 200, req.Code)
// 	}
// }

// func TestLDAPAuthenticationWithPassword(t *testing.T) {
// 	r := createGinAndController()

// 	req := helpers.PerformAuthRequest(r, "GET", "/api/v1/auth", "einstein", "password")

// 	// Test that the http status code is 200
// 	if req.Code != 200 {
// 		t.Errorf("Expected response code %d. Got %d\n", 200, req.Code)
// 	}
// }
