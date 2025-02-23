package user

import (
	"context"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/modasby/futeboxd-backend/core/pkg/errors"
	"github.com/modasby/futeboxd-backend/core/pkg/utils"
	"golang.org/x/crypto/bcrypt"
)

const (
	passwordCost int = 1
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
	ProfilePicture string
	Email          string
	Password       string
	FavoriteTeamID int64
}

func NewUser(ctx context.Context, username, name, bio, email, password string, favoriteTeamID int64) (*User, error) {
	user := &User{
		Username:       strings.ToLower(username),
		Name:           name,
		Bio:            bio,
		Email:          strings.ToLower(email),
		Password:       password,
		FavoriteTeamID: favoriteTeamID,
	}

	if err := user.Validate(ctx); err != nil {
		return nil, err
	}

	if err := user.HashPassword(); err != nil {
		return nil, err
	}

	return user, nil
}

func (u *User) UpdatePassword(ctx context.Context, newPassword string) error {
	if err := u.CheckPassword(ctx, newPassword); err == nil {
		httpErr := &errors.HTTPErr{
			Msg:        "A nova senha não pode ser igual a senha antiga",
			Code:       http.StatusBadRequest,
			StackTrace: errors.CaptureStackTrace(),
			Context:    "USER:DOMAIN:CHANGE_PASSWORD:SAME_PASSWORD",
			ErrorCode:  utils.GetTraceIDFromCtx(ctx),
			Timestamp:  time.Now().UTC(),
			Original:   err,
		}

		return httpErr
	}

	u.Password = newPassword

	if err := u.Validate(ctx); err != nil {
		return err
	}

	if err := u.HashPassword(); err != nil {
		return err
	}

	return nil
}

func (u *User) Validate(ctx context.Context) error {

	usernameValid, err := regexp.MatchString("^[a-zA-Z0-9_-]+$", u.Username)
	if err != nil {
		return err
	}

	if !usernameValid {
		httpErr := &errors.HTTPErr{
			Msg:        "username inválido",
			Code:       http.StatusBadRequest,
			StackTrace: errors.CaptureStackTrace(),
			Context:    "USER:DOMAIN:VALIDATE:INVALID_USERNAME",
			ErrorCode:  utils.GetTraceIDFromCtx(ctx),
			Timestamp:  time.Now().UTC(),
		}

		return httpErr
	}

	if bannedUsernames[u.Username] {
		httpErr := &errors.HTTPErr{
			Msg:        "nome de usuário não está disponível",
			Code:       http.StatusBadRequest,
			StackTrace: errors.CaptureStackTrace(),
			Context:    "USER:DOMAIN:VALIDATE:BANNED_USERNAME",
			ErrorCode:  utils.GetTraceIDFromCtx(ctx),
			Timestamp:  time.Now().UTC(),
		}

		return httpErr
	}

	matched, err := regexp.MatchString("^[\\w-.]+@([\\w-]+\\.)+[\\w-]{2,4}$", u.Email)
	if err != nil {
		return err
	}

	if !matched {
		httpErr := &errors.HTTPErr{
			Msg:        "Email inválido",
			Code:       http.StatusBadRequest,
			StackTrace: errors.CaptureStackTrace(),
			Context:    "USER:DOMAIN:VALIDATE:INVALID_EMAIL",
			ErrorCode:  utils.GetTraceIDFromCtx(ctx),
			Timestamp:  time.Now().UTC(),
		}

		return httpErr
	}

	if len(u.Password) <= 5 {
		httpErr := &errors.HTTPErr{
			Msg:        "senha deve ter pelo menos 5 caracteres",
			Code:       http.StatusBadRequest,
			StackTrace: errors.CaptureStackTrace(),
			Context:    "USER:DOMAIN:VALIDATE:INVALID_PASSWORD",
			ErrorCode:  utils.GetTraceIDFromCtx(ctx),
			Timestamp:  time.Now().UTC(),
		}

		return httpErr
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

func (u *User) CheckPassword(ctx context.Context, providedPassword string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(providedPassword)); err != nil {
		httpErr := &errors.HTTPErr{
			Msg:        "senha não confere",
			Code:       http.StatusBadRequest,
			StackTrace: errors.CaptureStackTrace(),
			Context:    "USER:DOMAIN:CHECK_PASSWORD:WRONG_PASSWORD",
			ErrorCode:  utils.GetTraceIDFromCtx(ctx),
			Timestamp:  time.Now().UTC(),
		}

		return httpErr
	}

	return nil
}
