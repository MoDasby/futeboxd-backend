package usecase

import (
	"log"

	"github.com/modasby/futeboxd-api/pkg/client/football"
	"github.com/modasby/futeboxd-api/services/users/internal/domain"
)

type FindUserByIdOrUsernameUseCase struct {
	userRepository domain.UserRepository
	footballClient football.Client
}

func NewFindUserByIdOrUsernameUseCase(
	userRepository domain.UserRepository,
	footballClient football.Client,
) *FindUserByIdOrUsernameUseCase {
	return &FindUserByIdOrUsernameUseCase{
		userRepository: userRepository,
		footballClient: footballClient,
	}
}

func (uc *FindUserByIdOrUsernameUseCase) Execute(identificator string) (*UserOutputDTO, error) {
	user, err := uc.userRepository.FindOneByIdOrUsername(identificator)
	if err != nil {
		return nil, err
	}

	favoriteTeam, err := uc.footballClient.GetTeam(user.FavoriteTeamID)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	output := UserOutputDTO{
		ID:           user.ID,
		Username:     user.Username,
		FavoriteTeam: favoriteTeam,
	}

	return &output, nil
}
