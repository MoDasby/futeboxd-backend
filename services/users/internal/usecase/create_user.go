package usecase

import (
	"github.com/modasby/futeboxd-api/pkg/client/football"
	"github.com/modasby/futeboxd-api/pkg/errors"
	"github.com/modasby/futeboxd-api/services/users/internal/domain"
)

type CreateUserInputDTO struct {
	Username       string `json:"username"`
	Email          string `json:"email"`
	Password       string `json:"password"`
	FavoriteTeamId int64  `json:"favorite_team_id"`
}

type UserOutputDTO struct {
	ID           string         `json:"id,omitempty"`
	Username     string         `json:"username"`
	Email        string         `json:"email,omitempty"`
	FavoriteTeam *football.Team `json:"favorite_team,omitempty"`
}

type CreateUserUseCase struct {
	userRepository domain.UserRepository
	footballClient football.Client
}

func NewCreateUserUseCase(
	userRepository domain.UserRepository,
	footballClient football.Client,
) *CreateUserUseCase {
	return &CreateUserUseCase{
		userRepository: userRepository,
		footballClient: footballClient,
	}
}

func (uc *CreateUserUseCase) Execute(input CreateUserInputDTO) (*UserOutputDTO, error) {
	user, err := domain.NewUser(input.Username, input.Email, input.Password, input.FavoriteTeamId)
	if err != nil {
		return nil, err
	}

	exists, err := uc.userRepository.Exists(user.Username, user.Email)
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, errors.NewErrBadRequest("já existe uma conta com esses dados")
	}

	var favoriteTeam *football.Team

	if input.FavoriteTeamId != 0 {
		favoriteTeam, err = uc.footballClient.GetTeam(input.FavoriteTeamId)
		if err != nil {
			return nil, err
		}
	}

	if err := uc.userRepository.AddUser(user); err != nil {
		return nil, err
	}

	output := UserOutputDTO{
		Username:     user.Username,
		Email:        user.Email,
		FavoriteTeam: favoriteTeam,
	}

	return &output, nil
}
