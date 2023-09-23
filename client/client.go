package client

import (
	"encoding/json"
	"github.com/candbright/client-auth/repo"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
)

const CodeSuccess = 0

type UserOperation interface {
	GetRegisterCode(phoneNumber string) (string, error)
	RegisterUser(phoneNumber string, code string, data repo.User) (repo.User, error)
	GetUserById(id string) (repo.User, error)
	GetUserByPhoneNumber(phoneNumber string) (repo.User, error)
	UpdateUserById(id string, data repo.User) error
	DeleteUserById(id string) error
}

type Client interface {
	AuthMiddleware(eng *gin.Engine, config *MiddlewareConfig) *gin.RouterGroup
	UserOperation
}

type client struct {
	endpoint string
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

func (client *client) RegisterUser(phoneNumber string, code string, data repo.User) (repo.User, error) {
	resp, err := client.client.R().
		SetQueryParams(map[string]string{
			"phone_number": phoneNumber,
			"code":         code,
		}).
		SetHeader("Accept", "application/json").
		SetBody(data).
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
	resp, err := client.client.R().
		SetQueryParams(map[string]string{
			"id": id,
		}).
		SetHeader("Accept", "application/json").
		Delete(client.endpoint + "/user")
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
