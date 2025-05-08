package domain

import (
	"time"

	"github.com/google/uuid"
)

type Person struct {
	ID          uuid.UUID
	Name        string
	Surname     string
	Age        	int
	Gender      string
	Nationality string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   time.Time
}

