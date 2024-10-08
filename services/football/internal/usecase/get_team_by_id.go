package usecase

import "github.com/modasby/futeboxd-api/services/football/internal/domain"

type GetTeamByIdUseCase struct {
	teamRepository domain.TeamRepository
}

func NewGetTeamByIdUseCase(teamRepository domain.TeamRepository) *GetTeamByIdUseCase {
	return &GetTeamByIdUseCase{
		teamRepository: teamRepository,
	}
}

func (uc *GetTeamByIdUseCase) Execute(teamID string) (*domain.Team, error) {
	return uc.teamRepository.FindTeamById(teamID)
}
