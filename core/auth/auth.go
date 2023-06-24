package auth

import (
	"encoding/json"
	"errors"
	"os"
)

type AuthType string

const (
	MASTER  AuthType = "master"
	USER    AuthType = "user"
	PREMIUM AuthType = "premium"
	BLOCKED AuthType = "blocked"
	ALL     AuthType = "*"
)

type User struct {
	Sender string
	Auth   AuthType
	Coins  int
	Level  int
	Hits   int
}

type Auther struct {
	SaveLocation string
	Database     map[string]map[string]User
}

type IAuther interface {
	LoadAuther() error
	SaveAuther() error
	ChatExists(string) bool
	SenderExists(string, string) bool
	ChangeAuthType(string, string, AuthType) error
	RemoveSender(string, string) error
	RegisterSender(User, string) error
	SumCoins(string, string, int, Calculate) error
	ChangeLevel(string, string, int, Calculate) error
	Hitting(string, string) error
	AuthPassed(string, string, AuthType) (bool, error)
	AllAuths() map[string]map[string]User
}

func NewAuther(p string) IAuther {
	var a = Auther{
		SaveLocation: p,
		Database:     map[string]map[string]User{},
	}

	a.LoadAuther()

	return &a
}

func (a *Auther) LoadAuther() error {
	var b []byte = []byte("{}")
	if bb, err := os.ReadFile(a.SaveLocation); err != nil {
		os.WriteFile(a.SaveLocation, b, os.ModeAppend)
	} else {
		b = bb
	}
	return json.Unmarshal(b, &a.Database)
}

func (a *Auther) SaveAuther() error {
	if b, err := json.MarshalIndent(a.Database, "", "  "); err != nil {
		return err
	} else {
		return os.WriteFile(a.SaveLocation, b, os.ModeAppend)
	}
}

func (a *Auther) ChatExists(chat string) bool {
	var _, ok = a.Database[chat]
	return ok
}

func (a *Auther) SenderExists(sender, chat string) bool {
	if a.ChatExists(chat) {
		return false
	} else {
		if _, okk := a.Database[chat][sender]; !okk {
			return false
		}
	}

	return true
}

func (a *Auther) ChangeAuthType(sender, chat string, t AuthType) error {

	if !a.SenderExists(sender, chat) {
		return errors.New("Auther Sender not exists")
	} else {
		if u, ok := a.Database[chat][sender]; ok {
			u.Auth = t
		}
		a.SaveAuther()
		return nil
	}
}

func (a *Auther) RemoveSender(sender, chat string) error {

	if !a.SenderExists(sender, chat) {
		return errors.New("Auther Sender not exists")
	} else {

		delete(a.Database[chat], sender)
		a.SaveAuther()
		return nil
	}
}

func (a *Auther) RegisterSender(u User, chat string) error {

	if !a.ChatExists(chat) {
		a.Database[chat] = map[string]User{}
	}

	if a.SenderExists(u.Sender, chat) {
		return errors.New("Auther Sender already exists")
	} else {
		a.Database[chat][u.Sender] = u
		a.SaveAuther()
		return nil
	}
}

type Calculate func(int, int) int

func Add(a, b int) int {
	return a + b
}

func Min(a, b int) int {
	return a - b
}

func (a *Auther) ChangeLevel(sender, chat string, level int, f Calculate) error {

	if !a.SenderExists(sender, chat) {
		return errors.New("Auther Sender not exists")
	} else {
		if u, ok := a.Database[chat][sender]; ok {
			u.Level = f(u.Level, level)
		}
	}

	return nil
}

func (a *Auther) SumCoins(sender, chat string, coin int, f Calculate) error {

	if !a.SenderExists(sender, chat) {
		return errors.New("Auther Sender not exists")
	} else {
		if u, ok := a.Database[chat][sender]; ok {
			u.Coins = f(u.Coins, coin)
		}
	}

	return nil
}

func (a *Auther) Hitting(sender, chat string) error {

	if !a.SenderExists(sender, chat) {
		return errors.New("Auther Sender not exists")
	} else {
		if u, ok := a.Database[chat][sender]; ok {
			u.Hits = Add(u.Hits, 1)
		}
	}

	return nil
}

func (a *Auther) AuthPassed(sender, chat string, t AuthType) (bool, error) {

	if !a.SenderExists(sender, chat) {
		return false, errors.New("Auther Sender not exists")
	} else {
		if u, ok := a.Database[chat][sender]; ok {
			return u.Auth == t || t == ALL, nil
		} else {
			return false, errors.New("Auther Sender not allowed")
		}
	}
}

func (a *Auther) AllAuths() map[string]map[string]User {

	return a.Database
}
