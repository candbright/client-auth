package client

import (
	"github.com/candbright/client-auth/repo"
	"net/url"
	"testing"
)

var phoneTest = "15000000000"
var userTest = repo.User{
	Username:    "test",
	PhoneNumber: phoneTest,
	Password:    "123456",
}

func TestClient_GetRegisterCode(t *testing.T) {
	resp, err := Ins().GetRegisterCode(phoneTest)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(resp)
}

func TestClient_RegisterOrLogin(t *testing.T) {
	resp, err := Ins().GetRegisterCode(phoneTest)
	if err != nil {
		t.Fatal(err)
	}
	user, err := Ins().RegisterOrLogin(phoneTest, resp)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(user)
}

func TestClient_DeleteUserById(t *testing.T) {
	code, err := Ins().GetRegisterCode(phoneTest)
	if err != nil {
		t.Fatal(err)
	}
	_, err = Ins().Login(url.Values{
		"phone_number": []string{phoneTest},
		"code":         []string{code},
	})
	if err != nil {
		t.Fatal(err)
	}
	user, err := Ins().GetUserByPhoneNumber(phoneTest)
	if err != nil {
		t.Fatal(err)
	}
	_, err = Ins().GetUserById(user.Id)
	if err != nil {
		t.Fatal(err)
	}
	err = Ins().DeleteUserById(user.Id)
	if err != nil {
		t.Fatal(err)
	}
	_, err = Ins().GetUserById(user.Id)
	t.Log(err)
}
