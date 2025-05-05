package dto

import (
	"time"
)

type CreatePersonRequest struct {
	Name    string `json:"name" binding:"required,min=2,max=50,alpha"`
	Surname string `json:"surname" binding:"required,min=2,max=50,alpha"`
}

type PersonResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Surname     string    `json:"surname"`
	Age         int       `json:"age"`
	Gender      string    `json:"gender"`
	Nationality string    `json:"nationality"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   time.Time `json:"deleted_at"`
}

type UpdatePersonRequest struct {
	Name        string `json:"name"`
	Surname     string `json:"surname"`
	Age         int    `json:"age"`
	Gender      string `json:"gender"`
	Nationality string `json:"nationality"`
}

type UpdatePersonResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Surname     string    `json:"surname"`
	Age         int       `json:"age"`
	Gender      string    `json:"gender"`
	Nationality string    `json:"nationality"`
	UpdateAt    time.Time `json:"updated_at"`
}
