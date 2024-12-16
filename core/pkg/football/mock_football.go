package football

import "github.com/stretchr/testify/mock"

type MockFootballClient struct {
	mock.Mock
}

func (m *MockFootballClient) GetMatchesMap(matchIDs []int64) (map[int64]Match, error) {
	return nil, nil
}

func (m *MockFootballClient) GetMatch(matchID int64) (*Match, error) {
	args := m.Called(matchID)

	match, _ := args.Get(0).(*Match)

	return match, args.Error(1)
}

func (m *MockFootballClient) GetTeam(teamID int64) (*Team, error) {
	args := m.Called(teamID)

	team, _ := args.Get(0).(*Team)

	return team, args.Error(1)
}

func (m *MockFootballClient) GetMatches(matchIDs []int64) ([]Match, error) {
	args := m.Called(matchIDs)

	matches, _ := args.Get(0).([]Match)

	return matches, nil
}
