package godoc

import (
	"errors"

	"github.com/beatlabs/gomodctl/internal"
	"github.com/go-resty/resty/v2"
)

type response struct {
	Results []struct {
		Name        string  `json:"name"`
		Path        string  `json:"path"`
		ImportCount int     `json:"import_count"`
		Stars       int     `json:"stars,omitempty"`
		Score       float64 `json:"score"`
		Synopsis    string  `json:"synopsis,omitempty"`
	} `json:"results"`
}

// Client is exported.
type Client struct {
	restClient *resty.Client
}

// NewClient is exported.
func NewClient() *Client {
	return &Client{restClient: resty.New()}
}

// Search is exported.
func (c *Client) Search(term string) ([]internal.SearchResult, error) {
	if term == "" {
		return nil, errors.New("empty term")
	}

	resp := &response{}

	_, err := c.restClient.R().
		SetQueryParams(map[string]string{
			"q": term,
		}).
		SetHeader("Accept", "application/json").
		SetResult(resp).
		Get("https://api.godoc.org/search")

	if err != nil {
		return nil, err
	}

	results := make([]internal.SearchResult, len(resp.Results))

	for i, result := range resp.Results {
		results[i] = internal.SearchResult{
			Name:        result.Name,
			Path:        result.Path,
			ImportCount: result.ImportCount,
			Stars:       result.Stars,
			Score:       result.Score,
			Synopsis:    result.Synopsis,
		}
	}

	return results, nil
}

// Info is exported.
func (c *Client) Info(path string) (string, error) {
	if path == "" {
		return "", errors.New("path is empty")
	}

	resp, err := c.restClient.R().
		SetHeader("Accept", "text/plain").
		Get("https://godoc.org/" + path)

	if err != nil {
		return "", err
	}

	response := resp.String()

	if response == "NOT FOUND" {
		return "", errors.New("not found")
	}

	return response, err
}
