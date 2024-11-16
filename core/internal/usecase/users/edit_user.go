package usecase

import (
	"github.com/modasby/futeboxd-api/services/core/internal/client/football"
	"github.com/modasby/futeboxd-api/services/core/internal/domain"
	"github.com/modasby/futeboxd-api/services/core/internal/errors"
)

type EditUserUseCase struct {
	userRepository domain.UserRepository
	footballClient football.Client
}

func NewEditUserUseCase(
	userRepository domain.UserRepository,
	footballClient football.Client,
) *EditUserUseCase {
	return &EditUserUseCase{
		userRepository: userRepository,
		footballClient: footballClient,
	}
}

type EditUserInputDTO struct {
	UserID         string
	Name           string `json:"name"`
	Bio            string `json:"bio"`
	Username       string `json:"username"`
	Email          string `json:"email,omitempty"`
	FavoriteTeamID int64  `json:"favorite_team_id"`
}

func (uc *EditUserUseCase) Execute(input EditUserInputDTO) error {
	user, err := uc.userRepository.FindOneByIdOrUsername(input.UserID)
	if err != nil {
		return err
	}

	if err := uc.updateEmail(user, input.Email); err != nil {
		return err
	}

	if err := uc.updateUsername(user, input.Username); err != nil {
		return err
	}

	if err := uc.updateFavoriteTeam(user, input.FavoriteTeamID); err != nil {
		return err
	}

	uc.updateBio(user, input.Bio)

	uc.updateName(user, input.Name)

	if err := user.Validate(); err != nil {
		return err
	}

	if err := uc.userRepository.Update(user); err != nil {
		return err
	}

	return nil
}

func (uc *EditUserUseCase) updateFavoriteTeam(user *domain.User, favoriteTeam int64) error {
	if favoriteTeam > 0 {
		team, err := uc.footballClient.GetTeam(favoriteTeam)
		if err != nil {
			return err
		}

		user.FavoriteTeamID = team.ID
	}

	return nil
}

func (uc *EditUserUseCase) updateEmail(user *domain.User, email string) error {
	if email != "" {
		emailExists, err := uc.userRepository.Exists("", email)
		if err != nil {
			return err
		}

		if emailExists && email != user.Email {
			return errors.NewHTTPErr(
				"esse email já existe",
				409,
				"USECASE:USER:EDIT:EMAIL_ALREADY_EXISTS",
			)
		}
		user.Email = email
	}

	return nil
}

func (uc *EditUserUseCase) updateUsername(user *domain.User, username string) error {
	if username != "" {
		usernameExists, err := uc.userRepository.Exists(username, "")
		if err != nil {
			return err
		}

		if usernameExists && username != user.Username {
			return errors.NewHTTPErr(
				"esse username já existe",
				409,
				"USECASE:USER:EDIT:USERNAME_ALREADY_EXISTS",
			)
		}

		user.Username = username
	}

	return nil
}

func (uc *EditUserUseCase) updateName(user *domain.User, name string) {
	if name != "" {
		user.Name = name
	}
}

func (uc *EditUserUseCase) updateBio(user *domain.User, bio string) {
	if bio != "" {
		user.Bio = bio
	}
}
