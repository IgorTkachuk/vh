package models

type User struct {
	Login    string `json:"login,omitempty"`
	Password string `json:"password,omitempty"`
}

type UserDto struct {
	Login    string `json:"login,omitempty"`
	Password string `json:"password,omitempty"`
}
