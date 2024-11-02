package domain

import (
	"math/rand"
	"time"
)

type UserId uint64

type User struct {
	Id           UserId
	Username     Username
	Password     Password
	CreatedUtc   time.Time
	BannedUtc    time.Time
	IsAdmin      bool
	LastLoginUtc time.Time
}

var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

func NewUser(username Username, password Password) User {
	return User{
		Id:           UserId(rng.Uint64()),
		Username:     username,
		Password:     password,
		CreatedUtc:   time.Now().UTC(),
		BannedUtc:    time.Time{},
		IsAdmin:      false,
		LastLoginUtc: time.Time{},
	}
}

func (u User) IsBanned() bool { return u.BannedUtc != time.Time{} }

func (u User) MatchPassword(plainText string) bool { return u.Password.MatchPassword(plainText) }

func (u *User) Ban() { u.BannedUtc = time.Now().UTC() }
