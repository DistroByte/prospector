package middleware

import (
	"fmt"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/go-ldap/ldap/v3"
)

func connectLdap(url string) (*ldap.Conn, error) {
	l, err := ldap.DialURL(url)
	if err != nil {
		fmt.Println("error connecting to ldap server:", err)
		return nil, err
	}

	return l, nil
}

func bindLdap(l *ldap.Conn, username string, password string) error {
	err := l.Bind(username, password)
	if err != nil {
		fmt.Println("error binding to ldap server:", err)
		return err
	}

	return nil
}

func AuthenticateLdap(username string, password string) (interface{}, error) {
	bindusername := "read-only-admin"
	bindpassword := "password"

	l, err := connectLdap("ldap://ldap.forumsys.com:389")
	if err != nil {
		return nil, jwt.ErrFailedAuthentication
	}

	defer l.Close()

	err = bindLdap(l, "cn="+bindusername+",dc=example,dc=com", bindpassword)
	if err != nil {
		return nil, jwt.ErrFailedAuthentication
	}

	searchRequest := ldap.NewSearchRequest(
		"dc=example,dc=com",
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(uid=%s)", username),
		[]string{"dn"},
		nil,
	)

	sr, err := l.Search(searchRequest)
	if err != nil {
		fmt.Println("error searching ldap server:", err)
		return nil, jwt.ErrFailedAuthentication
	}

	if len(sr.Entries) != 1 {
		fmt.Println("user does not exist")
		return nil, jwt.ErrFailedAuthentication
	}

	sr.PrettyPrint(2)
	userdn := sr.Entries[0].DN

	err = l.Bind(userdn, password)
	if err != nil {
		fmt.Printf("error binding to ldap server with user credentials %s\n", userdn)
	} else {
		return &User{
			Username: username,
			FistName: "Test",
			LastName: "User",
		}, nil
	}

	return nil, jwt.ErrFailedAuthentication
}
