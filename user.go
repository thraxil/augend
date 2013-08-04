package main

import (
	"code.google.com/p/go.crypto/bcrypt"
	"fmt"
	"github.com/tpjg/goriakpbc"
)

type User struct {
	Username   string
	Password   []byte
	Facts      riak.Many
	riak.Model `riak:"augend.user"`
}

func (u *User) Resolve(count int) (err error) {
	fmt.Println("resolve user")
	return nil
}

func NewUser(username, password string) *User {
	var user User
	err := riak.NewModel(username, &user)
	if err != nil {
		fmt.Println("could not create new user:", err)
		return nil
	}
	user.Username = username
	user.SetPassword(password)
	user.SaveAs(user.Username)
	return &user
}

//SetPassword takes a plaintext password and hashes it with bcrypt and sets the
//password field to the hash.
func (u *User) SetPassword(password string) {
	hpass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err) //this is a panic because bcrypt errors on invalid costs
	}
	u.Password = hpass
}

func (u User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword(u.Password, []byte(password))
	if err != nil {
		return false
	}
	return true
}

//Login validates and returns a user object if they exist in the database.
// func Login(ctx *Context, username, password string) (u *User, err error) {
//   err = ctx.C("users").Find(bson.M{"username": username}).One(&u)
//   if err != nil {
//       return
//   }

//   err = bcrypt.CompareHashAndPassword(u.Password, []byte(password))
//   if err != nil {
//       u = nil
//   }
//   return
// }
