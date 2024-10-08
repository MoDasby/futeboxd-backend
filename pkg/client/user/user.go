package user

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/modasby/futeboxd-api/pkg/errors"
)

type Client interface {
	GetAuthenticatedUser(token string) (*AuthenticateOutputDTO, error)
	GetUser(identificator string) (*UserOutputDTO, error)
	GetUsers(ids []string) ([]UserOutputDTO, error)
}

type userClient struct {
	baseUrl string
}

func NewClient(baseUrl string) Client {
	return &userClient{
		baseUrl: baseUrl,
	}
}

func (u *userClient) GetAuthenticatedUser(token string) (*AuthenticateOutputDTO, error) {
	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/auth/user", u.baseUrl), nil)

	req.Header.Add("Authorization", token)

	client := http.Client{}

	res, err := client.Do(req)
	if err != nil || res.StatusCode != 200 {
		return nil, errors.NewErrUnauthorized("token inválido")
	}
	defer res.Body.Close()

	var output AuthenticateOutputDTO

	if err := json.NewDecoder(res.Body).Decode(&output); err != nil {
		return nil, err
	}

	return &output, err
}

func (u *userClient) GetUser(identificator string) (*UserOutputDTO, error) {
	res, err := http.Get(fmt.Sprintf("%s/user/%s", u.baseUrl, identificator))
	if err != nil || res.StatusCode != 200 {
		return nil, errors.NewErrNotFound("usuário não encontrado")
	}
	defer res.Body.Close()

	var output UserOutputDTO

	if err := json.NewDecoder(res.Body).Decode(&output); err != nil {
		return nil, err
	}

	return &output, err
}

type FindUserBatchInput struct {
	IDs []string `json:"ids"`
}

func (u *userClient) GetUsers(ids []string) ([]UserOutputDTO, error) {
	body, _ := json.Marshal(FindUserBatchInput{
		IDs: ids,
	})

	payload := bytes.NewBuffer(body)

	res, err := http.Post(fmt.Sprintf("%s/batch/user", u.baseUrl), "application/json", payload)
	if err != nil {
		return nil, err
	}

	var output []UserOutputDTO

	if err := json.NewDecoder(res.Body).Decode(&output); err != nil {
		return nil, err
	}

	return output, nil
}
