package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"smo-gateway/internal/utils"
)

type GeneratorClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewGeneratorClient() *GeneratorClient {
	url := os.Getenv("GENERATOR_SERVICE_URL")
	if url == "" {
		url = "http://localhost:8081"
	}
	return &GeneratorClient{
		baseURL: url,
		httpClient: &http.Client{
			Timeout: 500 * time.Millisecond,
		},
	}
}

func (c *GeneratorClient) GetUniform(a, b float64) float64 {
	resp, err := c.httpClient.Get(fmt.Sprintf("%s/api/distr/uniform?a=%f&b=%f", c.baseURL, a, b))
	if err != nil || resp.StatusCode != http.StatusOK {
		return utils.Uniform(a, b) // fallback на локальную математику
	}
	defer resp.Body.Close()

	var res struct {
		Value float64 `json:"value"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return utils.Uniform(a, b)
	}
	return res.Value
}

func (c *GeneratorClient) GetExponential(lambda float64) float64 {
	resp, err := c.httpClient.Get(fmt.Sprintf("%s/api/distr/exponential?lambda=%f", c.baseURL, lambda))
	if err != nil || resp.StatusCode != http.StatusOK {
		return utils.Exponential(lambda) // fallback
	}
	defer resp.Body.Close()

	var res struct {
		Value float64 `json:"value"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return utils.Exponential(lambda)
	}
	return res.Value
}
