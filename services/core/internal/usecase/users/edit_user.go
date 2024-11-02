package usecase

import (
	"github.com/modasby/futeboxd-api/services/core/internal/client/football"
	"github.com/modasby/futeboxd-api/services/core/internal/domain"
	"github.com/modasby/futeboxd-api/services/core/internal/errors"
)

type EditUserUseCase struct {
	userRepository domain.UserRepository
	footballClient football.Client
}

func NewEditUserUseCase(
	userRepository domain.UserRepository,
	footballClient football.Client,
) *EditUserUseCase {
	return &EditUserUseCase{
		userRepository: userRepository,
		footballClient: footballClient,
	}
}

type EditUserInputDTO struct {
	UserID         string
	Username       string `json:"username"`
	Email          string `json:"email,omitempty"`
	FavoriteTeamID int64  `json:"favorite_team_id"`
}

func (uc *EditUserUseCase) Execute(input EditUserInputDTO) error {
	user, err := uc.userRepository.FindOneByIdOrUsername(input.UserID)
	if err != nil {
		return err
	}

	if input.Email != "" {
		emailExists, err := uc.userRepository.Exists("", input.Email)
		if err != nil {
			return err
		}

		if emailExists {
			return errors.NewHTTPErr(
				"esse email já existe",
				409,
				"USECASE:USER:EDIT:EMAIL_ALREADY_EXISTS",
			)
		}
		user.Email = input.Email
	}

	if input.Username != "" {
		usernameExists, err := uc.userRepository.Exists(input.Username, "")
		if err != nil {
			return err
		}

		if usernameExists {
			return errors.NewHTTPErr(
				"esse username já existe",
				409,
				"USECASE:USER:EDIT:USERNAME_ALREADY_EXISTS",
			)
		}

		user.Username = input.Username
	}

	if input.FavoriteTeamID > 0 {
		team, err := uc.footballClient.GetTeam(input.FavoriteTeamID)
		if err != nil {
			return err
		}

		user.FavoriteTeamID = team.ID
	}

	if err := user.Validate(); err != nil {
		return err
	}

	if err := uc.userRepository.Update(user); err != nil {
		return err
	}

	return nil
}
