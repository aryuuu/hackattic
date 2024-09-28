package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"time"
)

type HTTPClient struct {
	httpClient *http.Client
	baseURL    string
}

func NewGeoIPClient(baseURL string) *HTTPClient {
	return &HTTPClient{
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
		baseURL: baseURL,
	}
}

func (c *HTTPClient) GetProblem(ctx context.Context, path string, res interface{}) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+path, nil)
	if err != nil {
		log.Printf("failed to create req %w", err)
		return err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		log.Println("failed to get problem %w", err)
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("expected 200 status code, got %d", resp.StatusCode)
	}

	defer resp.Body.Close()

	c.BodyAsJSON(ctx, resp.Body, res)

	return nil
}

func (c *HTTPClient) PostSolution(ctx context.Context, path string, payload interface{}, res interface{}) error {
	buffer := bytes.NewBuffer([]byte{})
	if err := json.NewEncoder(buffer).Encode(payload); err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+path, buffer)
	if err != nil {
		log.Printf("failed to create req %w", err)
		return err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		log.Printf("failed to post solution %w", err)
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("expected 200 status code, got %d", resp.StatusCode)
	}

	resdump, err := httputil.DumpResponse(resp, true)
	if err != nil {
		log.Printf("failed to dump resp", err)
	}
	log.Println("res dump", string(resdump))

	defer resp.Body.Close()

	c.BodyAsJSON(ctx, resp.Body, res)

	return nil
}

func (c HTTPClient) BodyAsJSON(ctx context.Context, body io.ReadCloser, o interface{}) error {
	if body == nil {
		return nil
	}
	// read body as a JSON response
	raw, err := io.ReadAll(body)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(raw, o); err != nil {
		return err
	}

	return nil
}
