package usecase

import (
	"github.com/modasby/futeboxd-api/services/core/internal/domain"
)

type UpdatePasswordUsecase struct {
	userRepository domain.UserRepository
}

func NewUpdatePasswordUsecase(
	userRepository domain.UserRepository,
) *UpdatePasswordUsecase {
	return &UpdatePasswordUsecase{
		userRepository: userRepository,
	}
}

type UpdatePasswordInputDTO struct {
	User            *domain.User
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}

func (uc *UpdatePasswordUsecase) Execute(input UpdatePasswordInputDTO) error {
	user := input.User

	if err := user.UpdatePassword(input.CurrentPassword, input.NewPassword); err != nil {
		return err
	}

	return uc.userRepository.Update(user)
}
