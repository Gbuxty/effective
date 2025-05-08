package dto

import (
	"Effective/internal/domain"
	"errors"
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

func (req *UpdatePersonRequest) NewPerson(person *domain.Person) error {
	if person == nil {
		return errors.New("person is nil")
	}

	if req.Name != "" {
		person.Name = req.Name
	}

	if req.Surname != "" {
		person.Surname = req.Surname
	}

	if req.Age != 0 {
		person.Age = req.Age
	}

	if req.Gender != "" {
		person.Gender = req.Gender
	}

	if req.Nationality != "" {
		person.Nationality = req.Nationality
	}

	return nil
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
