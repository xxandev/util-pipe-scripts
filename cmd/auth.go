package main

import (
	"errors"
	"fmt"
	"util-pipe/internal/dbg"

	"github.com/go-ldap/ldap/v3"
)

var authorization sAuth

type sAuth struct {
	connLDAP *ldap.Conn
	IsLDAP   bool
	IsBase   bool
}

func (sa *sAuth) Init() error {
	sa.IsBase = len(config.BasicAuth.Login)+len(config.BasicAuth.Pass) > 0
	sa.IsLDAP = len(config.LDAP.URL)+len(config.LDAP.DN) > 0
	return nil
}

func (sa *sAuth) IsActive() bool {
	return sa.IsLDAP || sa.IsBase
}

func (sa *sAuth) CheckCredentials(login string, pass string, state bool) error {
	if !state {
		return errors.New("invalid credentials, auth empty")
	}
	if sa.IsBase && login == config.BasicAuth.Login && pass == config.BasicAuth.Pass {
		dbg.Log.Println("check credentials base auth - login:", login)
		return nil
	}
	if sa.IsLDAP {
		dbg.Log.Println("check credentials ldap auth - login:", login)
		if sa.connLDAP == nil {
			if err := sa.ConnLDAP(); err != nil {
				dbg.Log.Println("error ldap conn:", err)
				return err
			}
		}
		err := sa.connLDAP.Bind(fmt.Sprintf("uid=%v,%v", login, config.LDAP.DN), pass)
		if ldap.IsErrorWithCode(err, ldap.ErrorNetwork) {
			dbg.Log.Println("error ldap bind user:", err)
			dbg.Log.Println("error ldap conn close:", sa.connLDAP.Close())
			if err := sa.ConnLDAP(); err != nil {
				dbg.Log.Println("error ldap conn:", err)
				return err
			}
			return sa.connLDAP.Bind(fmt.Sprintf("uid=%v,%v", login, config.LDAP.DN), pass)
		}
		return err
	}
	return errors.New("invalid credentials")
}

func (sa *sAuth) ConnLDAP() error {
	var err error
	dbg.Log.Println("ldap conn:", config.LDAP.URL)
	sa.connLDAP, err = ldap.DialURL(config.LDAP.URL)
	return err
}
