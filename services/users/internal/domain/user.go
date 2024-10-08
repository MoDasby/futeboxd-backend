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
		return errors.NewErrUnauthorized("credencial ou senha incorretos")
	}
	return nil
}

func (u *User) Validate() error {
	if len(u.Username) <= 1 {
		return errors.NewErrBadRequest("username deve ter mais de 1 caracter")
	}

	matched, err := regexp.MatchString("^[\\w-\\.]+@([\\w-]+\\.)+[\\w-]{2,4}$", u.Email)
	if err != nil {
		return err
	}

	if !matched {
		return errors.NewErrBadRequest("email inválido")
	}

	if len(u.Password) <= 5 {
		return errors.NewErrBadRequest("senha deve ter pelo menos 5 caracteres")
	}

	return nil
}
