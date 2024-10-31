package usecase

import (
	"github.com/modasby/futeboxd-api/services/core/internal/domain"
)

type LogoutUsecase struct {
	sessionRepository domain.SessionRepository
}

func NewLogoutUsecase(
	sessionRepo domain.SessionRepository,
) *LogoutUsecase {
	return &LogoutUsecase{
		sessionRepository: sessionRepo,
	}
}

func (uc *LogoutUsecase) Execute(session *domain.Session) error {
	if err := uc.sessionRepository.Delete(session.ID); err != nil {
		return err
	}

	return nil
}
