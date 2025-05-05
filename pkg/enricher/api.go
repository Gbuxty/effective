package enricher

import (
	"Effective/config"
	"Effective/pkg/logger"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type Enricher struct {
	client *http.Client
	logger *logger.Logger
	cfg *config.Config
}

type agifyResponse struct {
	Age   int `json:"age"`
	Count int `json:"count"`
}

type genderizeResponse struct {
	Gender      string  `json:"gender"`
	Probability float64 `json:"probability"`
	Count       int     `json:"count"`
}

type nationalizeResponse struct {
	Country []nationalizeCountry `json:"country"`
}

type nationalizeCountry struct {
	CountryID   string  `json:"country_id"`
	Probability float64 `json:"probability"`
}

func New(logger *logger.Logger) *Enricher {
	return &Enricher{
		client: &http.Client{Timeout: 10 * time.Second},
		logger: logger,
	}
}

type EnrichedData struct {
	Age         int    `json:"age"`
	Gender      string `json:"gender"`
	Nationality string `json:"nationality"`
}

func (e *Enricher) Enrich(name string) (*EnrichedData, error) {
	age, err := e.getAge(name)
	if err != nil {
		return nil, fmt.Errorf("failed get age : %w", err)
	}

	gender, err := e.getGender(name)
	if err != nil {
		return nil, fmt.Errorf("failed get gender: %w", err)
	}

	national, err := e.getNationality(name)
	if err != nil {
		return nil, fmt.Errorf("failed get nationality: %w", err)
	}

	data := &EnrichedData{
		Age:         age,
		Gender:      gender,
		Nationality: national,
	}
	e.logger.Info("Enriched data success", zap.String("gender", gender), zap.Int("age", age), zap.String("national", national))

	return data, nil
}

// ссылки нужно в env
func (e *Enricher) getAge(name string) (int, error) {
	resp, err := e.client.Get(fmt.Sprintf("https://api.agify.io/?name=%s", name))
	if err != nil {
		return 0, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("unexpected status: %s", resp.Status)
	}

	result := &agifyResponse{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, fmt.Errorf("json decode failed: %w", err)
	}

	return result.Age, nil
}

func (e *Enricher) getGender(name string) (string, error) {
	resp, err := e.client.Get(fmt.Sprintf("https://api.genderize.io/?name=%s", name))
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)

	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status: %s", resp.Status)
	}

	result := &genderizeResponse{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("json decode failed: %w", err)
	}

	return result.Gender, nil
}

func (e *Enricher) getNationality(name string) (string, error) {
	resp, err := e.client.Get(fmt.Sprintf("https://api.nationalize.io/?name=%s", name))
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status: %s", resp.Status)
	}

	result := &nationalizeResponse{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("json decode failed: %w", err)
	}

	if len(result.Country) == 0 {
		result.Country = append(result.Country, nationalizeCountry{CountryID: "unknown"})
	}

	return result.Country[0].CountryID, nil
}
