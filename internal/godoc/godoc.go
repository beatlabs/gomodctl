package godoc

import (
	"context"
	"encoding/json"
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

type imports struct {
	Results []struct {
		Path string `json:"path"`
	} `json:"imports"`
}

type importers struct {
	Results []struct {
		Path string `json:"path"`
	} `json:"results"`
}

// Client is exported.
type Client struct {
	restClient *resty.Client
	ctx        context.Context
}

// NewClient is exported.
func NewClient(ctx context.Context) *Client {
	return &Client{restClient: resty.New(), ctx: ctx}
}

// Search is exported.
func (c *Client) Search(term string) ([]internal.SearchResult, error) {
	if term == "" {
		return nil, errors.New("empty term")
	}

	resp := &response{}

	_, err := c.restClient.R().
		SetContext(c.ctx).
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
		SetContext(c.ctx).
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

// Imports fetches imports of a package
func (c *Client) Imports(path string) ([]string, error) {
	if path == "" {
		return nil, errors.New("path is empty")
	}

	resp, err := c.restClient.R().
		SetContext(c.ctx).
		SetHeader("Accept", "application/json").
		Get("https://api.godoc.org/imports/" + path)

	if err != nil {
		return nil, err
	}

	imps := &imports{}
	if err := json.Unmarshal(resp.Body(), imps); err != nil {
		return nil, err
	}

	paths := []string{}
	for _, imp := range imps.Results {
		paths = append(paths, imp.Path)
	}
	return paths, nil
}

// Importers fetches importers of a package
func (c *Client) Importers(path string) ([]string, error) {
	if path == "" {
		return nil, errors.New("path is empty")
	}

	resp, err := c.restClient.R().
		SetContext(c.ctx).
		SetHeader("Accept", "application/json").
		Get("https://api.godoc.org/importers/" + path)

	if err != nil {
		return nil, err
	}

	imps := &importers{}
	if err := json.Unmarshal(resp.Body(), imps); err != nil {
		return nil, err
	}

	paths := []string{}
	for _, imp := range imps.Results {
		paths = append(paths, imp.Path)
	}
	return paths, nil
}
