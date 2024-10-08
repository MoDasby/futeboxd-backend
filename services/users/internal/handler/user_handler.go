package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/modasby/futeboxd-api/pkg/errors"
	"github.com/modasby/futeboxd-api/services/users/internal/usecase"
)

type UserHandler struct {
	createUserUseCase             *usecase.CreateUserUseCase
	loginUseCase                  *usecase.LoginUseCase
	findUserByIdOrUsernameUseCase *usecase.FindUserByIdOrUsernameUseCase
	authenticateUserUseCase       *usecase.AuthenticateUserUseCase
	findUserBatchUseCase          *usecase.FindUserBatchUseCase
}

func NewUserHandler(
	createUserUseCase *usecase.CreateUserUseCase,
	loginUseCase *usecase.LoginUseCase,
	findUserByIdOrUsernameUseCase *usecase.FindUserByIdOrUsernameUseCase,
	authenticateUserUseCase *usecase.AuthenticateUserUseCase,
	findUserBatchUseCase *usecase.FindUserBatchUseCase,
) *UserHandler {
	return &UserHandler{
		createUserUseCase:             createUserUseCase,
		loginUseCase:                  loginUseCase,
		findUserByIdOrUsernameUseCase: findUserByIdOrUsernameUseCase,
		authenticateUserUseCase:       authenticateUserUseCase,
		findUserBatchUseCase:          findUserBatchUseCase,
	}
}

func (h *UserHandler) AuthenticateUser(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")

	if token == "" {
		errors.HandleHttpError(w, errors.NewErrUnauthorized("token inválido"))

		return
	}

	output, err := h.authenticateUserUseCase.Execute(token)
	if err != nil {
		errors.HandleHttpError(w, err)

		return
	}

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(output)
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var input usecase.CreateUserInputDTO

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		errors.HandleHttpError(w, errors.NewErrBadRequest(err.Error()))

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

func (h *UserHandler) FindBatchByID(w http.ResponseWriter, r *http.Request) {
	var body usecase.FindUserBatchInput

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		errors.HandleHttpError(w, err)

		return
	}

	output, err := h.findUserBatchUseCase.Execute(body)
	if err != nil {
		errors.HandleHttpError(w, err)

		return
	}

	fmt.Println(output)

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(output)
}

func (h *UserHandler) FindByIdOrUsername(w http.ResponseWriter, r *http.Request) {
	identificator := r.PathValue("identificator")

	output, err := h.findUserByIdOrUsernameUseCase.Execute(identificator)
	if err != nil {
		errors.HandleHttpError(w, err)

		return
	}

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(output)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var loginInput usecase.LoginInput

	json.NewDecoder(r.Body).Decode(&loginInput)

	output, err := h.loginUseCase.Execute(&loginInput)
	if err != nil {
		errors.HandleHttpError(w, err)

		return
	}

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(output)
}
