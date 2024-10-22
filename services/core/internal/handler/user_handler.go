package handler

import (
	"encoding/json"
	"net/http"

	"github.com/modasby/futeboxd-api/pkg/errors"
	"github.com/modasby/futeboxd-api/services/core/internal/domain"
	"github.com/modasby/futeboxd-api/services/core/internal/dto"
	"github.com/modasby/futeboxd-api/services/core/internal/middleware"
	usecase "github.com/modasby/futeboxd-api/services/core/internal/usecase/users"
)

type UserHandler struct {
	createUserUseCase    *usecase.CreateUserUseCase
	findProfile          *usecase.FindProfile
	findUserBatchUseCase *usecase.FindUserBatchUseCase
	editUserUseCase      *usecase.EditUserUseCase
	loginUsecase         *usecase.LoginUseCase
	followUsecase        *usecase.FollowUserUseCase
	unfollowUsecase      *usecase.UnfollowUserUseCase
}

func NewUserHandler(
	createUserUseCase *usecase.CreateUserUseCase,
	findProfile *usecase.FindProfile,
	findUserBatchUseCase *usecase.FindUserBatchUseCase,
	editUserUseCase *usecase.EditUserUseCase,
	loginUsecase *usecase.LoginUseCase,
	followUsecase *usecase.FollowUserUseCase,
	unfollowUsecase *usecase.UnfollowUserUseCase,
) *UserHandler {
	return &UserHandler{
		createUserUseCase:    createUserUseCase,
		findProfile:          findProfile,
		findUserBatchUseCase: findUserBatchUseCase,
		editUserUseCase:      editUserUseCase,
		loginUsecase:         loginUsecase,
		followUsecase:        followUsecase,
		unfollowUsecase:      unfollowUsecase,
	}
}

func (h *UserHandler) RegisterRoutes(r *http.ServeMux, injectUser middleware.Middleware) {
	r.HandleFunc("/users/login", h.login)
	r.HandleFunc("/users", injectUser(h.getCurrentUser))
	r.HandleFunc("POST /users", h.createUser)
	r.HandleFunc("PUT /users", injectUser(h.editUser))
	r.HandleFunc("/profiles/{username}", injectUser(h.findByUsername))
	r.HandleFunc("POST /profiles/{username}/follow", injectUser(h.followUser))
	r.HandleFunc("DELETE /profiles/{username}/follow", injectUser(h.unfollowUser))
}

func (h *UserHandler) unfollowUser(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(middleware.UserKey).(*domain.User)
	if !ok || user == nil || user.Username == "anonymous" {
		httpErr := errors.NewHTTPErr("ocorreu um erro de autenticação, tente logar novamente", 401, "HANDLER:AUTHENTICATE_USER:INVALID_USER")
		errors.HandleHttpError(w, httpErr)

		return
	}

	followingUsername := r.PathValue("username")

	if err := h.unfollowUsecase.Execute(user.ID, followingUsername); err != nil {
		errors.HandleHttpError(w, err)

		return
	}

}

func (h *UserHandler) followUser(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(middleware.UserKey).(*domain.User)
	if !ok || user == nil || user.Username == "anonymous" {
		httpErr := errors.NewHTTPErr("ocorreu um erro de autenticação, tente logar novamente", 401, "HANDLER:AUTHENTICATE_USER:INVALID_USER")
		errors.HandleHttpError(w, httpErr)

		return
	}

	followingUsername := r.PathValue("username")

	if err := h.followUsecase.Execute(user.ID, followingUsername); err != nil {
		errors.HandleHttpError(w, err)

		return
	}
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
	if !ok || user == nil || user.Username == "anonymous" {
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
	if !ok || user == nil || user.Username == "anonymous" {
		httpErr := errors.NewHTTPErr("ocorreu um erro de autenticação, tente logar novamente", 401, "HANDLER:USER:EDIT:INVALID_SESSION")
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

	output, err := h.findProfile.Execute(requester, username)
	if err != nil {
		errors.HandleHttpError(w, err)

		return
	}

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(output)
}
