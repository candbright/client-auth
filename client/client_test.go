package client

import (
	"github.com/candbright/client-auth/repo"
	"testing"
)

func TestClient_GetRegisterCode(t *testing.T) {
	resp, err := Default().GetRegisterCode("15000000000")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(resp)
}

func TestClient_RegisterUser(t *testing.T) {
	phone := "15000000000"
	resp, err := Ins().GetRegisterCode(phone)
	if err != nil {
		t.Fatal(err)
	}
	userAdd, err := Ins().RegisterUser(phone, resp, repo.User{
		Username:    "test",
		PhoneNumber: phone,
		Password:    "123456",
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(userAdd)
}
