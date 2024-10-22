package usecase

import (
	"github.com/modasby/futeboxd-api/services/core/internal/domain"
	"github.com/modasby/futeboxd-api/services/core/internal/dto"
)

type EditUserUseCase struct {
	userRepository domain.UserRepository
}

func NewEditUserUseCase(userRepository domain.UserRepository) *EditUserUseCase {
	return &EditUserUseCase{
		userRepository: userRepository,
	}
}

type EditUserInputDTO struct {
	UserID string
	dto.UserInputDTO
}

func (uc *EditUserUseCase) Execute(input EditUserInputDTO) error {
	user, err := domain.NewUser(input.Username, input.Email, input.Password, input.FavoriteTeamID)
	if err != nil {
		return err
	}

	user.ID = input.UserID

	if err := uc.userRepository.EditUser(user); err != nil {
		return err
	}

	return nil
}
