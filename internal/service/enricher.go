package service

import (
	"Effective/config"
	"Effective/pkg/logger"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"
)

const (
	keyName = "name"
)

type Enricher struct {
	client *http.Client
	logger *logger.Logger
	cfg    *config.Config
}

func NewEnricher(logger *logger.Logger, cfg *config.Config) *Enricher {
	return &Enricher{
		client: &http.Client{Timeout: 5 * time.Second},
		logger: logger,
		cfg:    cfg,
	}
}

func (e *Enricher) GetAgeByName(ctx context.Context, name string) (int, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", e.cfg.APIUrl.AgifyUrl, nil)
	if err != nil {
		return 0, fmt.Errorf("request failed: %w", err)
	}

	query := req.URL.Query()
	query.Add(keyName, name)
	req.URL.RawQuery = query.Encode()
	e.logger.Debug("Request URL", zap.String("url", req.URL.String()))

	resp, err := e.client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("request failed: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("unexpected status: %s", resp.Status)
	}

	result := &apiAgeResponse{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, fmt.Errorf("json decode failed: %w", err)
	}
	e.logger.Debug("Response data from api", zap.Int("Age", result.Age))
	return result.Age, nil
}

func (e *Enricher) GetGenderByName(ctx context.Context, name string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", e.cfg.APIUrl.GenderizeUrl, nil)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}

	query := req.URL.Query()
	query.Add(keyName, name)
	req.URL.RawQuery = query.Encode()

	resp, err := e.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status: %s", resp.Status)
	}

	result := &apiGenderResponse{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("json decode failed: %w", err)
	}
	e.logger.Debug("Response data from api", zap.String("Gender", result.Gender))
	return result.Gender, nil
}

func (e *Enricher) GetNationalityByName(ctx context.Context, name string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", e.cfg.APIUrl.NationalizeUrl, nil)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}

	query := req.URL.Query()
	query.Add(keyName, name)
	req.URL.RawQuery = query.Encode()

	resp, err := e.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status: %s", resp.Status)
	}
	result := &apiNationalityResponse{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("json decode failed: %w", err)
	}

	if len(result.Country) == 0 {
		e.logger.Warn("No country data found, setting default value")
		result.Country = append(result.Country, nationalizeCountry{CountryID: "unknown"})
	}
	e.logger.Debug("Response data from api", zap.String("Nationality", result.Country[0].CountryID))
	return result.Country[0].CountryID, nil
}
