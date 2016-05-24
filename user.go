package main

import (
	"golang.org/x/crypto/bcrypt"
)

type user struct {
	Username string
	Password []byte
}

//SetPassword takes a plaintext password and hashes it with bcrypt and sets the
//password field to the hash.
func (u *user) SetPassword(password string) []byte {
	hpass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err) //this is a panic because bcrypt errors on invalid costs
	}
	u.Password = hpass
	return hpass
}

func (u user) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword(u.Password, []byte(password))
	if err != nil {
		return false
	}
	return true
}
