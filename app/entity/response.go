package entity

import "time"

//Response represents the all the response
type Response struct {
	Success bool        `json:"success" example:"true"`
	Message string      `json:"message" example:"created"`
	Data    interface{} `json:"data"`
}

type DataLogin struct {
	Token string `json:"token" example:"eyJhbGciOiJI...."`
}

type DataRegister struct {
	ID    uint   `json:"id" example:"1"`
	Email string `json:"email" example:"user@mail.com"`
	Uname string `json:"username" example:"user"`
	Age   int    `json:"age" example:"18"`
}

type DataPhoto struct {
	ID        uint        `json:"id" example:"1"`
	Title     string      `json:"title"`
	Caption   string      `json:"caption"`
	UserID    uint        `json:"id_user" example:"1"`
	Username  string      `json:"username"`
	Photo_URL string      `json:"photo_url"`
	CreatedAt *time.Time  `json:"created_at"`
	UpdatedAt *time.Time  `json:"updated_at"`
	Comment   interface{} `json:"comment"`
}

type DataComment struct {
	ID        uint       `json:"id" example:"1"`
	Message   string     `json:"message"`
	Username  string     `json:"username"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}
