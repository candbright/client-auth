package repo

type User struct {
	Id          string `json:"id"`
	PhoneNumber string `json:"phone_number"`
	Username    string `json:"user_name"`
	Password    string `json:"password"`
	CreateAt    string `json:"create_at"`
	UpdateAt    string `json:"update_at"`
}
