package domain

import (
	"regexp"
	"strings"

	"github.com/modasby/futeboxd-backend/core/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

const (
	passwordCost int = 10
)

var (
	bannedUsernames map[string]bool = map[string]bool{
		"partidas": true, "anonymous": true, "reviews": true,
	}
)

type User struct {
	ID             string
	Username       string
	Name           string
	Bio            string
	Email          string
	Password       string
	FavoriteTeamID int64
}

func NewUser(username, name, bio, email, password string, favoriteTeamID int64) (*User, error) {
	user := &User{
		Username:       strings.ToLower(username),
		Name:           name,
		Bio:            bio,
		Email:          strings.ToLower(email),
		Password:       password,
		FavoriteTeamID: favoriteTeamID,
	}

	if err := user.Validate(); err != nil {
		return nil, err
	}

	if err := user.HashPassword(); err != nil {
		return nil, err
	}

	return user, nil
}

func (u *User) UpdatePassword(newPassword string) error {
	if err := u.CheckPassword(newPassword); err == nil {
		return errors.NewHTTPErr(
			"a nova senha não pode ser igual a senha antiga",
			400,
			"DOMAIN:USER:CHANGE_PASSWORD:SAME_PASSWORD",
		)
	}

	u.Password = newPassword

	if err := u.Validate(); err != nil {
		return err
	}

	if err := u.HashPassword(); err != nil {
		return err
	}

	return nil
}

func (u *User) Validate() error {

	usernameValid, err := regexp.MatchString("^[a-zA-Z0-9_-]+$", u.Username)
	if err != nil {
		return err
	}

	if !usernameValid {
		return errors.NewHTTPErr(
			"username inválido",
			400,
			"DOMAIN:USER:VALIDATE:INVALID_USERNAME",
		)
	}

	if bannedUsernames[u.Username] {
		return errors.NewHTTPErr(
			"nome de usuário não está disponível",
			400,
			"DOMAIN:USER:VALIDADE:BANNED_USERNAME",
		)
	}

	matched, err := regexp.MatchString("^[\\w-.]+@([\\w-]+\\.)+[\\w-]{2,4}$", u.Email)
	if err != nil {
		return err
	}

	if !matched {
		return errors.NewHTTPErr("email inválido", 400, "DOMAIN:USER:VALIDATE:INVALID_EMAIL")
	}

	if len(u.Password) <= 5 {
		return errors.NewHTTPErr(
			"senha deve ter pelo menos 5 caracteres",
			400,
			"DOMAIN:USER:VALIDATE:INVALID_PASSWORD",
		)
	}

	if u.Name == "" {
		return errors.NewHTTPErr(
			"nome não pode estar vazio",
			400,
			"DOMAIN:USER:VALIDADE:EMPTY_NAME",
		)
	}

	return nil
}

func (u *User) HashPassword() error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(u.Password), passwordCost)
	if err != nil {
		return err
	}

	u.Password = string(bytes)

	return nil
}

func (u *User) CheckPassword(providedPassword string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(providedPassword)); err != nil {
		return errors.NewHTTPErr(
			"senha não confere",
			400,
			"DOMAIN:USER:CHECK_PASSWORD:WRONG_PASSWORD",
		)
	}

	return nil
}
