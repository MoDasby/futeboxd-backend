package usecase

import (
	"github.com/modasby/futeboxd-api/pkg/errors"
	"github.com/modasby/futeboxd-api/services/core/internal/domain"
	"github.com/modasby/futeboxd-api/services/core/internal/dto"
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

func (uc *FindUserBatchUseCase) Execute(input FindUserBatchInput) ([]dto.UserDTO, error) {
	if len(input.IDs) == 0 {
		return nil, errors.NewHTTPErr("corpo de requisição inválido", 400, "FIND_USER_BATCH:INVALID_BODY")
	}

	users, err := uc.userRepository.FindBatchByID(input.IDs)
	if err != nil {
		return nil, err
	}

	var output []dto.UserDTO

	for _, user := range users {
		output = append(output, dto.UserDTO{
			ID:       user.ID,
			Username: user.Username,
		})
	}

	return output, nil
}
