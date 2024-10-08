package usecase

import (
	"fmt"

	"github.com/modasby/futeboxd-api/services/users/internal/domain"
)

type FindUserBatchUseCase struct {
	userRepository domain.UserRepository
}

func NewFindUserBatchUseCase(userRepository domain.UserRepository) *FindUserBatchUseCase {
	return &FindUserBatchUseCase{
		userRepository: userRepository,
	}
}

type FindUserBatchInput struct {
	IDs []string `json:"ids"`
}

func (uc *FindUserBatchUseCase) Execute(input FindUserBatchInput) ([]UserOutputDTO, error) {

	users, err := uc.userRepository.FindBatchByID(input.IDs)
	if err != nil {
		return nil, err
	}

	fmt.Println(users)

	var output []UserOutputDTO

	for _, user := range users {
		output = append(output, UserOutputDTO{
			ID:       user.ID,
			Username: user.Username,
		})
	}

	return output, nil
}
