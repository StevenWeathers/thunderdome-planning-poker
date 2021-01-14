package main

import (
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/StevenWeathers/thunderdome-planning-poker/pkg/database"
	ldap "github.com/go-ldap/ldap/v3"
	"github.com/spf13/viper"
)

func (s *server) createCookie(warriorID string) *http.Cookie {
	encoded, err := s.cookie.Encode(s.config.SecureCookieName, warriorID)
	var NewCookie *http.Cookie

	if err == nil {
		NewCookie = &http.Cookie{
			Name:     s.config.SecureCookieName,
			Value:    encoded,
			Path:     s.config.PathPrefix + "/",
			HttpOnly: true,
			Domain:   s.config.AppDomain,
			MaxAge:   86400 * 30, // 30 days
			Secure:   s.config.SecureCookieFlag,
			SameSite: http.SameSiteStrictMode,
		}
	}
	return NewCookie
}

func (s *server) authWarriorDatabase(warriorEmail string, warriorPassword string) (*database.Warrior, error) {
	authedWarrior, err := s.database.AuthWarrior(warriorEmail, warriorPassword)
	if err != nil {
		log.Println("Failed authenticating user", warriorEmail)
	} else if authedWarrior == nil {
		log.Println("Unknown user", warriorEmail)
	}
	return authedWarrior, err
}

// Authenticate using LDAP and if warrior does not exist, automatically add warror as a verified warrior
func (s *server) authAndCreateWarriorLdap(warriorUsername string, warriorPassword string) (*database.Warrior, error) {
	var authedWarrior *database.Warrior
	l, err := ldap.DialURL(viper.GetString("auth.ldap.url"))
	if err != nil {
		log.Println("Failed connecting to ldap server at", viper.GetString("auth.ldap.url"))
		return authedWarrior, err
	}
	defer l.Close()
	if viper.GetBool("auth.ldap.use_tls") {
		err = l.StartTLS(&tls.Config{InsecureSkipVerify: true})
		if err != nil {
			log.Println("Failed securing ldap connection", err)
			return authedWarrior, err
		}
	}

	if viper.GetString("auth.ldap.bindname") != "" {
		err = l.Bind(viper.GetString("auth.ldap.bindname"), viper.GetString("auth.ldap.bindpass"))
		if err != nil {
			log.Println("Failed binding for authentication:", err)
			return authedWarrior, err
		}
	}

	searchRequest := ldap.NewSearchRequest(viper.GetString("auth.ldap.basedn"),
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf(viper.GetString("auth.ldap.filter"), warriorUsername),
		[]string{"dn", viper.GetString("auth.ldap.mail_attr"), viper.GetString("auth.ldap.cn_attr")},
		nil,
	)

	sr, err := l.Search(searchRequest)
	if err != nil {
		log.Println("Failed performing ldap search query for", warriorUsername, ":", err)
		return authedWarrior, err
	}

	if len(sr.Entries) != 1 {
		log.Println("User", warriorUsername, "does not exist or too many entries returned")
		return authedWarrior, errors.New("warrior not found")
	}

	userdn := sr.Entries[0].DN
	useremail := sr.Entries[0].GetAttributeValue(viper.GetString("auth.ldap.mail_attr"))
	usercn := sr.Entries[0].GetAttributeValue(viper.GetString("auth.ldap.cn_attr"))

	err = l.Bind(userdn, warriorPassword)
	if err != nil {
		log.Println("Failed authenticating user ", warriorUsername)
		return authedWarrior, err
	}

	authedWarrior, err = s.database.GetWarriorByEmail(useremail)
	if authedWarrior == nil {
		log.Println("Warrior", useremail, "does not exist in database, auto-recruit")
		newWarrior, verifyID, err := s.database.CreateWarriorCorporal(usercn, useremail, "", "")
		if err != nil {
			log.Println("Failed auto-creating new warrior", err)
			return authedWarrior, err
		}
		err = s.database.VerifyWarriorAccount(verifyID)
		if err != nil {
			log.Println("Failed verifying new warrior", err)
			return authedWarrior, err
		}
		authedWarrior = newWarrior
	}

	return authedWarrior, nil
}
