package model

import "fmt"

type User interface {
	Name() string
	Password() string
	Check(password string) bool
	fmt.Stringer
}

type UserImpl struct {
	Username string `json:"name" bson:"name"`
	Pass     string `json:"password" bson:"password"`
}

func (u UserImpl) Name() string {
	return u.Username
}

func (u UserImpl) Password() string {
	return u.Pass
}

func (u UserImpl) Check(password string) bool {
	return u.Password() == password
}

func (u UserImpl) String() string {
	return "{Username:" + u.Username + ";Pass:" + u.Pass + "}"
}
