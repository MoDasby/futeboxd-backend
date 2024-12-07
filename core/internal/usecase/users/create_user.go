package usecase

import (
	"github.com/modasby/futeboxd-api/services/core/internal/client/football"
	"github.com/modasby/futeboxd-api/services/core/internal/domain"
	"github.com/modasby/futeboxd-api/services/core/internal/dto"
	"github.com/modasby/futeboxd-api/services/core/internal/errors"
	"github.com/modasby/futeboxd-api/services/core/internal/json/null"
)

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

func (uc *CreateUserUseCase) Execute(input dto.UserInputDTO) (*dto.UserDTO, error) {
	user, err := domain.NewUser(
		input.Username,
		input.Name,
		input.Bio,
		input.Email,
		input.Password,
		input.FavoriteTeamID,
	)
	if err != nil {
		return nil, err
	}

	exists, err := uc.userRepository.Exists(user.Username, user.Email)
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, errors.NewHTTPErr(
			"já existe uma conta com esses dados",
			400,
			"USECASE:CREATE_USER:DUPLICATE",
		)
	}

	if input.FavoriteTeamID != 0 {
		// checa se o time existe
		_, err = uc.footballClient.GetTeam(input.FavoriteTeamID)
		if err != nil {
			return nil, err
		}
	}

	res, err := uc.userRepository.Create(user)
	if err != nil {
		return nil, err
	}

	output := dto.UserDTO{
		ID:       res.ID,
		Name:     null.NewString(user.Name),
		Username: user.Username,
		Email:    user.Email,
	}

	return &output, nil
}
