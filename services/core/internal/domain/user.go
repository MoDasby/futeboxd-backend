package domain

import (
	"regexp"

	"github.com/modasby/futeboxd-api/pkg/errors"
	"golang.org/x/crypto/bcrypt"
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

func (u *User) HashPassword() error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(u.Password), 10)
	if err != nil {
		return err
	}

	u.Password = string(bytes)

	return nil
}

func (user *User) CheckPassword(providedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(providedPassword))
	if err != nil {
		return errors.NewHTTPErr(
			"credencial ou senha incorretos",
			401,
			"DOMAIN:USER:CHECK_PASSWORD:WRONG_PASSWORD",
		)
	}
	return nil
}

func (u *User) Validate() error {
	if len(u.Username) <= 1 {
		return errors.NewHTTPErr(
			"username deve ter mais de 1 caracter",
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
