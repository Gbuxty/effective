package service

import (
	"Effective/internal/domain"
	"Effective/internal/transport/http/handler/dto"
	"Effective/pkg/logger"
	"context"
	"fmt"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

type PersonService struct {
	repo     PersonRepository
	logger   *logger.Logger
	enricher EnricherService
}

type PersonRepository interface {
	SavePerson(ctx context.Context, person *domain.Person) (uuid.UUID, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Person, error)
	DeleteByID(ctx context.Context, id uuid.UUID) (bool, error)
	UpdatePerson(ctx context.Context, person *domain.Person) error
	GetPersonFilter(ctx context.Context, person *domain.PersonFilter) (*[]domain.Person, error)
}

type EnricherService interface {
	GetAgeByName(ctx context.Context, name string) (int, error)
	GetGenderByName(ctx context.Context, name string) (string, error)
	GetNationalityByName(ctx context.Context, name string) (string, error)
}

func NewPersonService(repo PersonRepository, logger *logger.Logger, enricher EnricherService) *PersonService {
	return &PersonService{
		repo:     repo,
		logger:   logger,
		enricher: enricher,
	}
}

func (s *PersonService) CreatePerson(ctx context.Context, req *dto.CreatePersonRequest) (uuid.UUID, error) {
	dataEnrichment := make(chan EnrichmentData, 1)

	g, gctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		enrichedAge, err := s.enricher.GetAgeByName(gctx, req.Name)
		if err != nil {
			return fmt.Errorf("failed to enrich age: %w", err)
		}
		dataEnrichment <- EnrichmentData{Type: "age", Value: enrichedAge}
		return nil
	})

	g.Go(func() error {
		enrichedGender, err := s.enricher.GetGenderByName(gctx, req.Name)
		if err != nil {
			return fmt.Errorf("failed to enrich gender: %w", err)
		}
		dataEnrichment <- EnrichmentData{Type: "gender", Value: enrichedGender}
		return nil
	})

	g.Go(func() error {
		enrichedNationality, err := s.enricher.GetNationalityByName(gctx, req.Name)
		if err != nil {
			return fmt.Errorf("failed to enrich nationality: %w", err)
		}
		dataEnrichment <- EnrichmentData{Type: "nationality", Value: enrichedNationality}
		return nil
	})

	go func() {
		err := g.Wait()
		if err != nil {
			s.logger.Error("Failed to enrich data", zap.Error(err))
		}
		close(dataEnrichment)
	}()

	person := &domain.Person{
		Name:    req.Name,
		Surname: req.Surname,
	}

	for data := range dataEnrichment {
		switch data.Type {
		case "age":
			if age, ok := data.Value.(int); ok {
				person.Age = age
			}
		case "gender":
			if gender, ok := data.Value.(string); ok {
				person.Gender = gender
			}
		case "nationality":
			if nationality, ok := data.Value.(string); ok {
				person.Nationality = nationality
			}

		}
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

func (s *PersonService) UpdatePerson(ctx context.Context, id uuid.UUID, req *dto.UpdatePersonRequest) error {
	person, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get person: %w", err)
	}

	if err := req.NewPerson(person); err != nil {
		return fmt.Errorf("failed to map person:%w", err)
	}

	if err := s.repo.UpdatePerson(ctx, person); err != nil {
		return fmt.Errorf("failed to update person:%w", err)
	}

	return nil
}

func (s *PersonService) GetPersonWithFilter(ctx context.Context, filter *dto.Filter) (*[]domain.Person, error) {
	personFilter := &domain.PersonFilter{
		Name:        filter.Name,
		Surname:     filter.Surname,
		MinAge:      filter.MinAge,
		MaxAge:      filter.MaxAge,
		Gender:      filter.Gender,
		Nationality: filter.Nationality,
		Page:        filter.Page,
		Size:        filter.Size,
	}

	filterPerson, err := s.repo.GetPersonFilter(ctx, personFilter)
	if err != nil {
		return nil, fmt.Errorf("failed to get person with filter:%w", err)
	}

	return filterPerson, nil
}
