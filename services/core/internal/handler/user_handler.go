package handler

import (
	"encoding/json"
	"net/http"

	"github.com/modasby/futeboxd-api/services/core/internal/domain"
	"github.com/modasby/futeboxd-api/services/core/internal/dto"
	"github.com/modasby/futeboxd-api/services/core/internal/errors"
	"github.com/modasby/futeboxd-api/services/core/internal/middleware"
	usecase "github.com/modasby/futeboxd-api/services/core/internal/usecase/users"
)

type UserHandler struct {
	createUserUseCase     *usecase.CreateUserUseCase
	findProfileUsecase    *usecase.FindProfileUsecase
	editUserUseCase       *usecase.EditUserUseCase
	loginUsecase          *usecase.LoginUseCase
	logoutUsecase         *usecase.LogoutUsecase
	updatePasswordUsecase *usecase.UpdatePasswordUsecase
	toggleFollowUsecase   *usecase.ToggleFollowUsecase
}

func NewUserHandler(
	createUserUseCase *usecase.CreateUserUseCase,
	findProfileUsecase *usecase.FindProfileUsecase,
	editUserUseCase *usecase.EditUserUseCase,
	loginUsecase *usecase.LoginUseCase,
	logoutUsecase *usecase.LogoutUsecase,
	updatePasswordUsecase *usecase.UpdatePasswordUsecase,
	toggleFollowUsecase *usecase.ToggleFollowUsecase,
) *UserHandler {
	return &UserHandler{
		createUserUseCase:     createUserUseCase,
		findProfileUsecase:    findProfileUsecase,
		editUserUseCase:       editUserUseCase,
		loginUsecase:          loginUsecase,
		logoutUsecase:         logoutUsecase,
		updatePasswordUsecase: updatePasswordUsecase,
		toggleFollowUsecase:   toggleFollowUsecase,
	}
}

func (h *UserHandler) RegisterRoutes(r *http.ServeMux, injectUser middleware.Middleware) {
	r.HandleFunc("POST /users/login", h.login)
	r.HandleFunc("POST /users/logout", injectUser(h.logout, false))
	r.HandleFunc("GET /users", injectUser(h.getCurrentUser, false))
	r.HandleFunc("POST /users", h.createUser)
	r.HandleFunc("PATCH /users", injectUser(h.editUser, false))
	r.HandleFunc("PUT /users/password", injectUser(h.updatePassword, false))
	r.HandleFunc("GET /profiles/{username}", injectUser(h.findByUsername, true))
	r.HandleFunc("POST /profiles/{username}/follow", injectUser(h.toggleFollow, false))
}

func (h *UserHandler) updatePassword(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(middleware.UserKey).(*domain.User)
	if !ok || user == nil {
		httpErr := errors.NewHTTPErr("ocorreu um erro de autenticação, tente logar novamente", 401, "HANDLER:AUTHENTICATE_USER:INVALID_USER")
		errors.HandleHttpError(w, httpErr)

		return
	}

	var input usecase.UpdatePasswordInputDTO

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		err = errors.NewHTTPErr(
			"corpo de requisição inválido",
			400,
			"HANDLER:UPDATE_PASSWORD:INVALID_BODY",
		)
	}

	input.User = user

	if err := h.updatePasswordUsecase.Execute(input); err != nil {
		errors.HandleHttpError(w, err)

		return
	}
}

func (h *UserHandler) logout(w http.ResponseWriter, r *http.Request) {
	session, ok := r.Context().Value(middleware.SessionKey).(*domain.Session)
	if !ok {
		httpErr := errors.NewHTTPErr("ocorreu um erro de autenticação, tente logar novamente", 401, "HANDLER:AUTHENTICATE_USER:INVALID_USER")
		errors.HandleHttpError(w, httpErr)

		return
	}

	if err := h.logoutUsecase.Execute(session); err != nil {
		errors.HandleHttpError(w, err)

		return
	}

}

func (h *UserHandler) toggleFollow(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(middleware.UserKey).(*domain.User)
	if !ok || user == nil {
		httpErr := errors.NewHTTPErr("ocorreu um erro de autenticação, tente logar novamente", 401, "HANDLER:AUTHENTICATE_USER:INVALID_USER")
		errors.HandleHttpError(w, httpErr)

		return
	}

	followingUsername := r.PathValue("username")

	output, err := h.toggleFollowUsecase.Execute(user.ID, followingUsername)
	if err != nil {
		errors.HandleHttpError(w, err)

		return
	}

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(output)
}

func (h *UserHandler) login(w http.ResponseWriter, r *http.Request) {
	var loginInput usecase.LoginInput

	json.NewDecoder(r.Body).Decode(&loginInput)

	output, err := h.loginUsecase.Execute(&loginInput)
	if err != nil {
		errors.HandleHttpError(w, err)

		return
	}

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(output)
}

func (h *UserHandler) getCurrentUser(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(middleware.UserKey).(*domain.User)
	if !ok || user == nil {
		httpErr := errors.NewHTTPErr("ocorreu um erro de autenticação, tente logar novamente", 401, "HANDLER:AUTHENTICATE_USER:INVALID_USER")
		errors.HandleHttpError(w, httpErr)

		return
	}

	type res struct {
		ID       string `json:"id"`
		Email    string `json:"email"`
		Username string `json:"username"`
	}

	response := res{
		ID:       user.ID,
		Email:    user.Email,
		Username: user.Username,
	}

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *UserHandler) createUser(w http.ResponseWriter, r *http.Request) {
	var input dto.UserInputDTO

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		errors.HandleHttpError(w, errors.NewHTTPErr("corpo de requisição inválido", 400, "HANDLER:USER:CREATE:INVALID_BODY"))

		return
	}

	output, err := h.createUserUseCase.Execute(input)
	if err != nil {
		errors.HandleHttpError(w, err)

		return
	}

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(output)
}

func (h *UserHandler) editUser(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(middleware.UserKey).(*domain.User)
	if !ok || user == nil {
		httpErr := errors.NewHTTPErr("ocorreu um erro de autenticação, tente logar novamente", 401, "HANDLER:AUTHENTICATE_USER:INVALID_USER")
		errors.HandleHttpError(w, httpErr)

		return
	}

	var input usecase.EditUserInputDTO

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		errors.HandleHttpError(w, errors.NewHTTPErr("corpo de requisição inválido", 400, "HANDLER:USER:EDIT:INVALID_BODY"))

		return
	}

	input.UserID = user.ID

	err := h.editUserUseCase.Execute(input)
	if err != nil {
		errors.HandleHttpError(w, err)

		return
	}

	w.WriteHeader(200)
}

func (h *UserHandler) findByUsername(w http.ResponseWriter, r *http.Request) {
	username := r.PathValue("username")

	requester, ok := r.Context().Value(middleware.UserKey).(*domain.User)
	if !ok || requester == nil {
		httpErr := errors.NewHTTPErr("ocorreu um erro de autenticação, tente logar novamente", 401, "HANDLER:USER:EDIT:INVALID_SESSION")
		errors.HandleHttpError(w, httpErr)

		return
	}

	output, err := h.findProfileUsecase.Execute(requester, username)
	if err != nil {
		errors.HandleHttpError(w, err)

		return
	}

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(output)
}
