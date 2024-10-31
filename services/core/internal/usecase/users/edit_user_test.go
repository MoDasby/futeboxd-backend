package usecase

import (
	"testing"
)

func TestEditUser(t *testing.T) {
	/* mockRepo := new(repository.MockUserRepo)
	mockFootballClient := new(football.MockFootballClient)

	uc := NewEditUserUseCase(mockRepo, mockFootballClient)

	input := EditUserInputDTO{
		UserID:         "uuid",
		Username:       "modasby",
		Email:          "modasby@email.com",
		FavoriteTeamID: 7632,
	}

	mockRepo.On("Update", mock.AnythingOfType("*domain.User")).Return(nil)

	err := uc.Execute(input)

	mockRepo.AssertCalled(t, "Update", mock.MatchedBy(func(user *domain.User) bool {
		return user.ID == "uuid" && user.Username == "modasby" &&
			user.Email == "modasby@email.com" && user.FavoriteTeamID == 7632
	}))

	assert.Nil(t, err)

	mockRepo.AssertExpectations(t) */
}
