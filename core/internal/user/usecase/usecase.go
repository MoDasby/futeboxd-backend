package usecase

import (
	"context"
	"net/http"
	"time"

	"github.com/modasby/futeboxd-backend/core/config"
	"github.com/modasby/futeboxd-backend/core/internal/user"
	"github.com/modasby/futeboxd-backend/core/internal/user/dto"
	"github.com/modasby/futeboxd-backend/core/pkg/email"
	"github.com/modasby/futeboxd-backend/core/pkg/email/templates"
	"github.com/modasby/futeboxd-backend/core/pkg/errors"
	"github.com/modasby/futeboxd-backend/core/pkg/football"
	"github.com/modasby/futeboxd-backend/core/pkg/json/null"
	"github.com/modasby/futeboxd-backend/core/pkg/token"
	"github.com/modasby/futeboxd-backend/core/pkg/utils"
)

type usersUsecases struct {
	userRepo       user.Repository
	footballClient football.Client
	emailClient    email.Client
	frontendConfig config.Frontend
}

func NewUsersUsecases(
	userRepo user.Repository,
	footballClient football.Client,
	emailClient email.Client,
	frontendConfig config.Frontend,
) user.Usecases {
	return &usersUsecases{
		userRepo:       userRepo,
		footballClient: footballClient,
		emailClient:    emailClient,
		frontendConfig: frontendConfig,
	}
}

func (uc *usersUsecases) CreateUser(ctx context.Context, input *dto.UserInput) error {
	user, err := user.NewUser(
		ctx,
		input.Username,
		input.Name,
		input.Bio,
		input.Email,
		input.Password,
		input.FavoriteTeamID,
	)
	if err != nil {
		return err
	}

	user.ProfilePicture = "https://api.dicebear.com/9.x/initials/png?seed=" + user.Username

	exists, err := uc.userRepo.Exists(ctx, user.Username, user.Email)
	if err != nil {
		return err
	}

	if exists {
		httpErr := &errors.HTTPErr{
			Msg:        "Já existe uma conta com esses dados",
			Code:       http.StatusConflict,
			StackTrace: errors.CaptureStackTrace(),
			Context:    "USER:USECASE:CREATE_USER:DUPLICATE",
			ErrorCode:  utils.GetTraceIDFromCtx(ctx),
			Timestamp:  time.Now().UTC(),
			Original:   err,
		}
		return httpErr
	}

	if input.FavoriteTeamID != 0 {
		_, err = uc.footballClient.GetTeam(ctx, input.FavoriteTeamID)
		if err != nil {
			return err
		}
	}

	if err := uc.userRepo.Create(ctx, user); err != nil {
		return err
	}

	return nil
}

func (uc *usersUsecases) EditUser(ctx context.Context, input *dto.EditUser) error {
	session, err := utils.GetSessionFromCtx(ctx)
	if err != nil {
		return err
	}

	user, err := uc.userRepo.FindOneByIdOrUsername(ctx, session.UserID)
	if err != nil {
		return err
	}

	if input.Email != "" {
		emailExists, err := uc.userRepo.Exists(ctx, "", input.Email)
		if err != nil {
			return err
		}

		if emailExists && input.Email != user.Email {
			httpErr := &errors.HTTPErr{
				Msg:        "Esse email já existe",
				Code:       http.StatusConflict,
				StackTrace: errors.CaptureStackTrace(),
				Context:    "USER:USECASE:EDIT:EMAIL_ALREADY_EXISTS",
				ErrorCode:  utils.GetTraceIDFromCtx(ctx),
				Timestamp:  time.Now().UTC(),
				Original:   err,
			}
			return httpErr
		}
		user.Email = input.Email
	}

	if input.Username != "" {
		usernameExists, err := uc.userRepo.Exists(ctx, input.Username, "")
		if err != nil {
			return err
		}

		if usernameExists && input.Username != user.Username {
			httpErr := &errors.HTTPErr{
				Msg:        "Esse username já existe",
				Code:       http.StatusConflict,
				Context:    "USER:USECASE:EDIT:USERNAME_ALREADY_EXISTS",
				StackTrace: errors.CaptureStackTrace(),
				ErrorCode:  utils.GetTraceIDFromCtx(ctx),
				Timestamp:  time.Now().UTC(),
				Original:   nil,
			}

			return httpErr
		}

		user.Username = input.Username
	}

	if input.FavoriteTeamID > 0 {
		team, err := uc.footballClient.GetTeam(ctx, input.FavoriteTeamID)
		if err != nil {
			return err
		}

		user.FavoriteTeamID = team.ID
	}

	if input.Bio != "" {
		user.Bio = input.Bio
	}

	if input.Name != "" {
		user.Name = input.Name
	}

	if err := user.Validate(ctx); err != nil {
		return err
	}

	if err := uc.userRepo.Update(ctx, user); err != nil {
		return err
	}

	return nil
}

func (uc *usersUsecases) ChangePassword(ctx context.Context, input *dto.ChangePassword) error {
	session, err := utils.GetSessionFromCtx(ctx)
	if err != nil {
		return err
	}

	user, err := uc.userRepo.FindOneByIdOrUsername(ctx, session.UserID)
	if err != nil {
		return err
	}

	if err := user.CheckPassword(ctx, input.CurrentPassword); err != nil {
		httpErr := &errors.HTTPErr{
			Msg:        "Senha incorreta",
			Code:       http.StatusUnauthorized,
			StackTrace: errors.CaptureStackTrace(),
			Context:    "USER:USECASE:CHANGE_PASSWORD:CHECK_PASSWORD:WRONG_PASSWORD",
			ErrorCode:  utils.GetTraceIDFromCtx(ctx),
			Timestamp:  time.Now().UTC(),
			Original:   err,
		}

		return httpErr
	}

	if err := user.UpdatePassword(ctx, input.NewPassword); err != nil {
		return err
	}

	return uc.userRepo.Update(ctx, user)
}

func (uc *usersUsecases) GetCurrentUser(ctx context.Context) (*dto.User, error) {
	session, err := utils.GetSessionFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	user, err := uc.userRepo.FindOneByIdOrUsername(ctx, session.UserID)
	if err != nil {
		return nil, err
	}

	var team *football.Team

	if user.FavoriteTeamID > 0 {
		team, err = uc.footballClient.GetTeam(ctx, user.FavoriteTeamID)
		if err != nil {
			return nil, err
		}
	}

	return &dto.User{
		ID:             user.ID,
		Name:           null.String(user.Name),
		Bio:            null.NewString(user.Bio),
		FavoriteTeam:   team,
		Username:       user.Username,
		Email:          user.Email,
		ProfilePicture: user.ProfilePicture,
	}, nil
}

func (uc *usersUsecases) RecoverPassword(ctx context.Context, input *dto.RecoverPassword) error {
	token, err := token.Generate()
	if err != nil {
		return err
	}

	requester, err := uc.userRepo.FindOneByCredential(ctx, input.Email)
	if err != nil {
		return err
	}

	recover := user.Recover{
		Token:  token,
		UserID: requester.ID,
	}

	if err := uc.userRepo.SaveRecoverToken(ctx, &recover); err != nil {
		return err
	}

	template := templates.NewRecoverPasswordTemplate(
		requester.Username,
		uc.frontendConfig.URL+"/recuperar/"+token,
	)

	mailData := email.Metadata{
		To:        requester.Email,
		Subject:   "Recuperação de senha",
		Body:      template,
		LocalPart: "support",
	}

	if err := uc.emailClient.SendMail(ctx, mailData); err != nil {
		return err
	}

	return nil
}

func (uc *usersUsecases) ResetPassword(ctx context.Context, input *dto.ResetPassword) error {
	recover, err := uc.userRepo.CheckRecoverToken(ctx, input.Token)
	if err != nil {
		return err
	}

	if recover.IsExpired() {
		httpErr := &errors.HTTPErr{
			Msg:        "O token de recuperação informado está expirado ou não existe. Solicite uma nova recuperação de senha.",
			Code:       http.StatusForbidden,
			StackTrace: errors.CaptureStackTrace(),
			Context:    "USER:USECASE:RESET_PASSWORD:EXPIRED_TOKEN",
			ErrorCode:  utils.GetTraceIDFromCtx(ctx),
			Timestamp:  time.Now().UTC(),
			Original:   nil,
		}
		return httpErr
	}

	user, err := uc.userRepo.FindOneByIdOrUsername(ctx, recover.UserID)
	if err != nil {
		return err
	}

	if err := user.UpdatePassword(ctx, input.Password); err != nil {
		return err
	}

	if err := uc.userRepo.Update(ctx, user); err != nil {
		return err
	}

	return nil
}
