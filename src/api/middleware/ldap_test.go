package middleware

import (
	"testing"

	"github.com/go-ldap/ldap/v3"
)

func TestLdapConnection(t *testing.T) {
	_, err := ldapConnection()
	if err != nil {
		t.Errorf("Test failed, expected nil, got %s", err)
	}
}

func TestBindLdap(t *testing.T) {
	l, _ := ldapConnection()

	tcs := []struct {
		name     string
		username string
		password string
		result   string
	}{
		{
			name:     "bind with correct credentials",
			username: "cn=read-only-admin,dc=example,dc=com",
			password: "password",
			result:   "",
		},
		{
			name:     "bind with incorrect credentials",
			username: "cn=read-only-admin,dc=example,dc=com",
			password: "wrongpassword",
			result:   "LDAP Result Code 49 \"Invalid Credentials\": ",
		},
	}

	for _, tc := range tcs {
		err := bindLdap(l, tc.username, tc.password)
		if err != nil {
			if err.Error() != tc.result {
				t.Errorf("Test failed, expected `%s`, got `%s`", tc.result, err)
			}
		}
	}
}

func TestAuthenticateLdap(t *testing.T) {
	tcs := []struct {
		name     string
		username string
		password string
		bindUser string
		bindPass string
		result   string
	}{
		{
			name:     "authenticate with correct credentials",
			username: "tesla",
			password: "password",
			bindUser: "read-only-admin",
			bindPass: "password",
			result:   "",
		},
		{
			name:     "authenticate with incorrect credentials",
			username: "tesla",
			password: "wrongpassword",
			bindUser: "read-only-admin",
			bindPass: "password",
			result:   "incorrect Username or Password",
		},
		{
			name:     "authenticate with incorrect bind credentials",
			username: "tesla",
			password: "password",
			bindUser: "read-only-admin",
			bindPass: "wrongpassword",
			result:   "incorrect Username or Password",
		},
		{
			name:     "authenticate with incorrect username",
			username: "wrongusername",
			password: "password",
			bindUser: "read-only-admin",
			bindPass: "password",
			result:   "incorrect Username or Password",
		},
	}

	for _, tc := range tcs {
		_, err := AuthenticateLdap(tc.username, tc.password, tc.bindUser, tc.bindPass)
		if err != nil {
			if err.Error() != tc.result {
				t.Errorf("Test failed, expected `%s`, got `%s`", tc.result, err)
			}
		}
	}
}

func ldapConnection() (*ldap.Conn, error) {
	url := "ldap://ldap.forumsys.com:389"
	conn, err := connectLdap(url)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
