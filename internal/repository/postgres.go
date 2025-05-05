package repository

import (
	"Effective/internal/domain"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PersonRepository struct {
	db *pgxpool.Pool
}

var (
	ErrUserNotFound = errors.New("user not found")
)

func NewPersonRepository(db *pgxpool.Pool) *PersonRepository {
	return &PersonRepository{db: db}
}

func (r *PersonRepository) SavePerson(ctx context.Context, person *domain.Person) (uuid.UUID, error) {
	var id uuid.UUID

	query := `
		INSERT INTO persons (
			name,
			surname,
			age,
			gender,
			nationality,
			created_at,
			updated_at
		) VALUES (
			$1, $2, $3, $4, $5, NOW(), NOW()
		)
		RETURNING id`

	err := r.db.QueryRow(
		ctx,
		query,
		person.Name,
		person.Surname,
		person.Age,
		person.Gender,
		person.Nationality,
	).Scan(&id)

	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to create user: %w", err)
	}

	return id, nil
}

func (r *PersonRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Person, error) {
	person := domain.Person{}

	query := `
			SELECT 
				id,
				name,
				surname,
				gender,
				nationality,
				created_at,
				updated_at
			FROM persons
			WHERE id=$1
			`
	err := r.db.QueryRow(
		ctx,
		query,
		id,
	).Scan(&person.ID,
		&person.Name,
		&person.Surname,
		&person.Gender,
		&person.Nationality,
		&person.CreatedAt,
		&person.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get person by id: %w", err)
	}
	return &person, nil
}

func (r *PersonRepository) DeleteByID(ctx context.Context, id uuid.UUID) (bool, error) {
	query := `UPDATE persons SET deleted_at = NOW() WHERE id = $1`

	_, err := r.db.Exec(ctx, query, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, ErrUserNotFound
		}
		return false, fmt.Errorf("failed to delete person by id: %w", err)
	}
	return true, nil
}

func (r *PersonRepository) UpdatePerson(ctx context.Context, person *domain.Person) (*domain.Person, error) {
	updatePerson := domain.Person{}

	query := `UPDATE persons 
				SET 
					name = $1,
					surname = $2,
					age = $3,
					gender = $4,
					nationality = $5,
					updated_at = NOW()
				WHERE id = $6
				RETURNING id, name, surname, age, gender, nationality, updated_at`

	err := r.db.QueryRow(
		ctx,
		query,
		person.Name,
		person.Surname,
		person.Age,
		person.Gender,
		person.Nationality,
		person.ID,
	).Scan(&updatePerson.ID,
		&updatePerson.Name,
		&updatePerson.Surname,
		&updatePerson.Age,
		&updatePerson.Gender,
		&updatePerson.Nationality,
		&updatePerson.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to update person: %w", err)
	}
	return &updatePerson, nil
}

func (r *PersonRepository) GetPersonFilter(ctx context.Context, person *domain.PersonFilter) (*domain.Person, error) {

	filterPerson := domain.Person{}
	query:=`SELECT id,name,surname,age`


	return nil, nil
}
