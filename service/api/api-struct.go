package api

import (
	"github.com/wippy00/wasa-text/service/database"
)

type User struct {
	Id    int64  `json:"id"`
	Name  string `json:"name"`
	Photo []byte `json:"photo"`
}

func NewUser(user database.User) User {
	return User{
		Id:    user.Id,
		Name:  user.Name,
		Photo: user.Photo,
	}
}
