package usecase

import "github.com/modasby/futeboxd-api/services/football/internal/domain"

type ListTeamsUseCase struct {
	teamRepository domain.TeamRepository
}

func NewListTeamsUseCase(teamRepository domain.TeamRepository) *ListTeamsUseCase {
	return &ListTeamsUseCase{
		teamRepository: teamRepository,
	}
}

func (uc *ListTeamsUseCase) Execute(name string) ([]*domain.Team, error) {

	teams, err := uc.teamRepository.FindAll(name)
	if err != nil {
		return nil, err
	}

	return teams, nil
}
