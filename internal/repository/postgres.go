package repository

import (
	"Effective/internal/domain"
	"context"
	"errors"
	"fmt"
	"log"

	sq "github.com/Masterminds/squirrel"
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

func (r *PersonRepository) UpdatePerson(ctx context.Context, person *domain.Person) error {
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

	_, err := r.db.Exec(
		ctx,
		query,
		person.Name,
		person.Surname,
		person.Age,
		person.Gender,
		person.Nationality,
		person.ID,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrUserNotFound
		}
		return fmt.Errorf("failed to update person: %w", err)
	}
	return nil
}

func (r *PersonRepository) GetPersonFilter(ctx context.Context, person *domain.PersonFilter) (*[]domain.Person, error) {
	query := sq.Select("id,name,surname,age,gender,nationality, created_at, updated_at").From("persons").PlaceholderFormat(sq.Dollar)

	if person.Name != nil {
		query = query.Where(sq.Eq{"name": *person.Name})
	}

	if person.Surname != nil {
		query = query.Where(sq.Eq{"surname": *person.Surname})
	}

	if person.MinAge != nil && person.MaxAge != nil {
		query = query.Where(sq.And{
			sq.GtOrEq{"age": *person.MinAge},
			sq.LtOrEq{"age": *person.MaxAge},
		})
	} else if person.MinAge != nil {
		query = query.Where(sq.GtOrEq{"age": *person.MinAge})
	} else if person.MaxAge != nil {
		query = query.Where(sq.LtOrEq{"age": *person.MaxAge})
	}

	if person.MaxAge != nil {
		query = query.Where(sq.LtOrEq{"age": *person.MaxAge})
	}
	if person.Gender != nil {
		query = query.Where(sq.Eq{"gender": *person.Gender})
	}
	if person.Nationality != nil {
		query = query.Where(sq.Eq{"nationality": *person.Nationality})
	}

	if person.Page <= 0 {
		person.Page = 1
	}
	if person.Size <= 0 {
		person.Size = 10
	}

	offset := (person.Page - 1) * person.Size
	query = query.Offset(uint64(offset)).Limit(uint64(person.Size))

	q, values, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}
	log.Printf("SQL: %s, Args: %v", q, values)

	filterPerson := make([]domain.Person, 0)

	rows, err := r.db.Query(ctx, q, values...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get person by filter: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var pers domain.Person
		if err := rows.Scan(
			&pers.ID,
			&pers.Name,
			&pers.Surname,
			&pers.Age,
			&pers.Gender,
			&pers.Nationality,
			&pers.CreatedAt,
			&pers.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan person: %w", err)
		}
		filterPerson = append(filterPerson, pers)
	}

	if len(filterPerson) == 0 {
		return nil, fmt.Errorf("no person found with the given filter")
	}
	return &filterPerson, nil
}
