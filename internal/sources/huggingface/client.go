package huggingface

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/pablodiegoo/dataprofile/internal/domain"
)

// Client queries the public Hugging Face Hub REST API.
type Client struct {
	HTTPClient *http.Client
	BaseURL    string
}

// NewClient returns a Hugging Face Hub client.
func NewClient() *Client {
	return &Client{
		HTTPClient: &http.Client{Timeout: 12 * time.Second},
		BaseURL:    "https://huggingface.co",
	}
}

// Summary holds aggregated counts across models, datasets, and spaces.
type Summary struct {
	Models       int
	Datasets     int
	Spaces       int
	Downloads30d int
	Likes        int
}

type hfItem struct {
	ID        string `json:"id"`
	Downloads int    `json:"downloads"`
	Likes     int    `json:"likes"`
}

// Collect queries models, datasets, and spaces published by the configured Hugging Face user.
func (c *Client) Collect(ctx context.Context, cfg domain.HuggingFaceSourceConfig) (Summary, error) {
	var out Summary
	username := strings.TrimSpace(cfg.Username)
	if !cfg.Enabled || username == "" {
		return out, nil
	}

	models, err := c.fetchEndpoint(ctx, fmt.Sprintf("%s/api/models?author=%s", c.BaseURL, username))
	if err == nil {
		out.Models = len(models)
		for _, m := range models {
			out.Downloads30d += m.Downloads
			out.Likes += m.Likes
		}
	}

	datasets, err := c.fetchEndpoint(ctx, fmt.Sprintf("%s/api/datasets?author=%s", c.BaseURL, username))
	if err == nil {
		out.Datasets = len(datasets)
		for _, d := range datasets {
			out.Downloads30d += d.Downloads
			out.Likes += d.Likes
		}
	}

	spaces, err := c.fetchEndpoint(ctx, fmt.Sprintf("%s/api/spaces?author=%s", c.BaseURL, username))
	if err == nil {
		out.Spaces = len(spaces)
		for _, s := range spaces {
			out.Likes += s.Likes
		}
	}

	return out, nil
}

func (c *Client) fetchEndpoint(ctx context.Context, url string) ([]hfItem, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("hf status %d", resp.StatusCode)
	}

	var list []hfItem
	if err := json.NewDecoder(resp.Body).Decode(&list); err != nil {
		return nil, err
	}
	return list, nil
}
