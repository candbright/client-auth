package client

import (
	"encoding/json"
	"github.com/candbright/client-auth/repo"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
	"net/url"
)

const CodeSuccess = 0

type UserOperation interface {
	GetRegisterCode(phoneNumber string) (string, error)
	RegisterOrLogin(phoneNumber, code string) (repo.User, error)
	GetUserById(id string) (repo.User, error)
	GetUserByPhoneNumber(phoneNumber string) (repo.User, error)
	UpdateUserById(id string, data repo.User) error
	DeleteUserById(id string) error
}

type LoginOperation interface {
	Login(params url.Values) (string, error)
	Logout(user repo.User) error
	RefreshToken(token string) error
}

type Client interface {
	AuthMiddleware(eng gin.IRouter, config *MiddlewareConfig) gin.IRouter
	UserOperation
	LoginOperation
}

type client struct {
	endpoint string
	token    string
	client   *resty.Client
}

func (client *client) GetRegisterCode(phoneNumber string) (string, error) {
	resp, err := client.client.R().
		SetQueryParam("phone_number", phoneNumber).
		SetHeader("Accept", "application/json").
		Get(client.endpoint + "/register")
	if err != nil {
		return "", errors.WithStack(err)
	}
	var result repo.Result[string]
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return "", errors.WithStack(err)
	}
	if result.Code != CodeSuccess {
		return "", errors.New(result.Message)
	}
	return result.Data, nil
}

func (client *client) RegisterOrLogin(phoneNumber, code string) (repo.User, error) {
	resp, err := client.client.R().
		SetHeader("Accept", "application/json").
		SetQueryParam("phone_number", phoneNumber).
		SetQueryParam("code", code).
		Post(client.endpoint + "/register")
	if err != nil {
		return repo.User{}, errors.WithStack(err)
	}
	var result repo.Result[repo.User]
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return repo.User{}, errors.WithStack(err)
	}
	if result.Code != CodeSuccess {
		return repo.User{}, errors.New(result.Message)
	}
	return result.Data, nil
}

func (client *client) GetUserById(id string) (repo.User, error) {
	resp, err := client.client.R().
		SetQueryParam("id", id).
		SetHeader("Accept", "application/json").
		SetHeader("Authorization", "Bearer "+client.token).
		Get(client.endpoint + "/user")
	if err != nil {
		return repo.User{}, errors.WithStack(err)
	}
	var result repo.Result[repo.User]
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return repo.User{}, errors.WithStack(err)
	}
	if result.Code != CodeSuccess {
		return repo.User{}, errors.New(result.Message)
	}
	return result.Data, nil
}

func (client *client) GetUserByPhoneNumber(phoneNumber string) (repo.User, error) {
	resp, err := client.client.R().
		SetQueryParam("phone_number", phoneNumber).
		SetHeader("Accept", "application/json").
		SetHeader("Authorization", "Bearer "+client.token).
		Get(client.endpoint + "/user")
	if err != nil {
		return repo.User{}, errors.WithStack(err)
	}
	var result repo.Result[repo.User]
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return repo.User{}, errors.WithStack(err)
	}
	if result.Code != CodeSuccess {
		return repo.User{}, errors.New(result.Message)
	}
	return result.Data, nil
}

func (client *client) UpdateUserById(id string, data repo.User) error {
	resp, err := client.client.R().
		SetQueryParams(map[string]string{
			"id": id,
		}).
		SetHeader("Accept", "application/json").
		SetHeader("Authorization", "Bearer "+client.token).
		SetBody(data).
		Put(client.endpoint + "/user")
	if err != nil {
		return errors.WithStack(err)
	}
	var result repo.Result[repo.User]
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return errors.WithStack(err)
	}
	if result.Code != CodeSuccess {
		return errors.New(result.Message)
	}
	return nil
}

func (client *client) DeleteUserById(id string) error {
	_, err := client.client.R().
		SetQueryParams(map[string]string{
			"id": id,
		}).
		SetHeader("Accept", "application/json").
		SetHeader("Authorization", "Bearer "+client.token).
		Delete(client.endpoint + "/user")
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (client *client) Login(params url.Values) (string, error) {
	resp, err := client.client.R().
		SetHeader("Accept", "application/json").
		SetQueryParamsFromValues(params).
		Post(client.endpoint + "/login")
	if err != nil {
		return "", errors.WithStack(err)
	}
	var result map[string]interface{}
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return "", errors.WithStack(err)
	}
	token, ok := result["token"].(string)
	if !ok {
		return "", errors.New("token is nil")
	}
	client.token = token
	return token, nil
}

func (client *client) Logout(user repo.User) error {
	_, err := client.client.R().
		SetHeader("Accept", "application/json").
		SetBody(user).
		Post(client.endpoint + "/logout")
	if err != nil {
		return errors.WithStack(err)
	}
	client.token = ""
	return nil
}

func (client *client) RefreshToken(token string) error {
	resp, err := client.client.R().
		SetHeader("Accept", "application/json").
		SetAuthToken(token).
		Get(client.endpoint + "/refresh_token")
	if err != nil {
		return errors.WithStack(err)
	}
	var result map[string]interface{}
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return errors.WithStack(err)
	}
	tokenRefresh, ok := result["token"].(string)
	if !ok {
		return errors.New("token is nil")
	}
	client.token = tokenRefresh
	return nil
}
