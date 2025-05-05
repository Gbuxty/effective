package service

import (
	"Effective/internal/domain"
	"Effective/internal/repository"
	"Effective/internal/transport/http/handler/dto"

	"Effective/pkg/enricher"
	"Effective/pkg/logger"
	"context"
	"fmt"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type PersonService struct {
	repo     *repository.PersonRepository
	enricher *enricher.Enricher
	logger   *logger.Logger
}

func NewPersonService(repo *repository.PersonRepository, enricher *enricher.Enricher, logger *logger.Logger) *PersonService {
	return &PersonService{
		repo:     repo,
		enricher: enricher,
		logger:   logger,
	}
}

func (s *PersonService) CreatePerson(ctx context.Context, req *dto.CreatePersonRequest) (uuid.UUID, error) {
	s.logger.Info("Enriching data", zap.String("name", req.Name))
	enriched, err := s.enricher.Enrich(req.Name)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to enrich data: %w", err)
	}

	person := &domain.Person{
		Name:        req.Name,
		Surname:     req.Surname,
		Age:         enriched.Age,
		Gender:      enriched.Gender,
		Nationality: enriched.Nationality,
	}

	id, err := s.repo.SavePerson(ctx, person)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to save person: %w", err)
	}

	return id, nil
}

func (s *PersonService) DeletePerson(ctx context.Context, id uuid.UUID) (bool, error) {
	_, err := s.repo.DeleteByID(ctx, id)
	if err != nil {
		return false, fmt.Errorf("failed to delete person:%w", err)
	}

	return true, nil
}

func (s *PersonService) UpdatePerson(ctx context.Context, id uuid.UUID, req *dto.UpdatePersonRequest) (*dto.UpdatePersonResponse, error) {
	person, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get person:%w", err)
	}
	// вот как мне эту хуйню красиво написать?
	person.Name = req.Name
	person.Surname = req.Surname
	person.Age = req.Age
	person.Gender = req.Gender
	person.Nationality = req.Nationality
	//затираются поля которые не передаются в запросе
	/* if req.Name != "" {
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
	} */

	changePerson, err := s.repo.UpdatePerson(ctx, person)
	if err != nil {
		return nil, fmt.Errorf("failed to upadate person in service :%w", err)
	}

	resp := &dto.UpdatePersonResponse{
		ID:          changePerson.ID.String(),
		Name:        changePerson.Name,
		Surname:     changePerson.Surname,
		Age:         changePerson.Age,
		Gender:      changePerson.Gender,
		Nationality: changePerson.Nationality,
		UpdateAt:    changePerson.UpdatedAt,
	}

	return resp, nil
}

func (s *PersonService) GetPersonWithFilter(ctx context.Context, filter *dto.Filter) (*dto.PersonResponse, error) {
	personFilter := &domain.PersonFilter{
		Name:        &filter.Name,
		Surname:     &filter.Surname,
		MinAge:      &filter.MinAge,
		MaxAge:      &filter.MaxAge,
		Gender:      &filter.Gender,
		Nationality: &filter.Nationality,
	}

	filterPerson, err := s.repo.GetPersonFilter(ctx, personFilter)
	if err != nil {
		return nil, fmt.Errorf("failed to get person with filter:%w", err)
	}

	resp := &dto.PersonResponse{
		ID:      filterPerson.ID.String(),
		Name:    filterPerson.Name,
		Surname: filterPerson.Surname,
		Age:     filterPerson.Age,
		Gender:  filter.Gender,
	}

	return resp, nil
}
