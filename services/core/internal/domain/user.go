package domain

import (
	"regexp"
	"strings"

	"github.com/modasby/futeboxd-api/services/core/internal/errors"
	"golang.org/x/crypto/bcrypt"
)

const (
	passwordCost int = 10
)

type User struct {
	ID             string
	Username       string
	Email          string
	Password       string
	FavoriteTeamID int64
}

func NewAnonymousUser() *User {
	return &User{
		ID:       "anonymous",
		Username: "anonymous",
	}
}

func NewUser(username string, email string, password string, favoriteTeamID int64) (*User, error) {
	user := &User{
		Username:       username,
		Email:          email,
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

func (u *User) UpdatePassword(currentPassword, newPassword string) error {

	if err := u.CheckPassword(currentPassword); err != nil {
		return errors.NewHTTPErr(
			"senhas não conferem",
			401,
			"DOMAIN:USER:UPDATE_PASSWORD:WRONG_PASSWORD",
		)
	}

	if err := u.CheckPassword(newPassword); err == nil {
		return errors.NewHTTPErr(
			"a nova senha não pode ser igual a senha antiga",
			400,
			"DOMAIN:USER:UPDATE_PASSWORD:SAME_PASSWORD",
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
	if strings.Contains(u.Username, " ") {
		return errors.NewHTTPErr(
			"o username não pode conter espaços",
			400,
			"DOMAIN:USER:VALIDADE:WHITE_SPACE_IN_USERNAME",
		)
	}

	matched, err := regexp.MatchString("^[\\w-\\.]+@([\\w-]+\\.)+[\\w-]{2,4}$", u.Email)
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
			"username deve ter mais de 1 caracter",
			400,
			"DOMAIN:USER:VALIDATE:INVALID_USERNAME",
		)
	}

	if strings.Contains(u.Username, " ") {
		return errors.NewHTTPErr(
			"username são pode conter espaços",
			400,
			"DOMAIN:USER:VALIDATE:INVALID_USERNAME",
		)
	}

	matched, err := regexp.MatchString("^[\\w-\\.]+@([\\w-]+\\.)+[\\w-]{2,4}$", u.Email)
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

	return nil
}
