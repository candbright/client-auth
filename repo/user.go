package repo

import (
	"github.com/google/uuid"
	"time"
)

func NewUser() *User {
	random, _ := uuid.NewRandom()
	return &User{
		Id: random.String(),
	}
}

type User struct {
	Id          string    `json:"id" gorm:"id;primaryKey"`
	PhoneNumber string    `json:"phone_number" gorm:"phone_number"`
	Username    string    `json:"username" gorm:"username"`
	Password    string    `json:"password" gorm:"password"`
	CreateAt    time.Time `json:"create_at" gorm:"create_at"`
	UpdateAt    time.Time `json:"update_at" gorm:"update_at"`
}

func (user User) TableName() string {
	return "user"
}
