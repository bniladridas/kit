package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const defaultTimeout = 30 * time.Second

var (
	ErrUnauthorized       = errors.New("unauthorized")
	ErrForbidden          = errors.New("forbidden")
	ErrNotFound           = errors.New("not found")
	ErrConflict           = errors.New("conflict")
	ErrUnprocessableEntity = errors.New("unprocessable entity")
	ErrRateLimited        = errors.New("rate limited")
)

type Client struct {
	BaseURL    string
	Token      string
	HTTPClient *http.Client
}

type RequestOption func(*http.Request)

func NewClient(baseURL, token string) *Client {
	if !strings.HasSuffix(baseURL, "/") {
		baseURL += "/"
	}
	return &Client{
		BaseURL: baseURL,
		Token:   token,
		HTTPClient: &http.Client{
			Timeout: defaultTimeout,
		},
	}
}

func (c *Client) Get(ctx context.Context, path string, opts ...RequestOption) (*Response, error) {
	return c.doRequest(ctx, http.MethodGet, path, nil, opts...)
}

func (c *Client) Post(ctx context.Context, path string, body interface{}, opts ...RequestOption) (*Response, error) {
	return c.doRequest(ctx, http.MethodPost, path, body, opts...)
}

func (c *Client) Patch(ctx context.Context, path string, body interface{}, opts ...RequestOption) (*Response, error) {
	return c.doRequest(ctx, http.MethodPatch, path, body, opts...)
}

func (c *Client) Delete(ctx context.Context, path string, opts ...RequestOption) (*Response, error) {
	return c.doRequest(ctx, http.MethodDelete, path, nil, opts...)
}

func (c *Client) doRequest(ctx context.Context, method, path string, body interface{}, opts ...RequestOption) (*Response, error) {
	reqURL := c.BaseURL + strings.TrimPrefix(path, "/")

	var bodyReader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(data)
	}

	req, err := http.NewRequestWithContext(ctx, method, reqURL, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "kit/0.1")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}

	if bodyReader != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	for _, opt := range opts {
		opt(req)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}

	apiResp := &Response{
		StatusCode: resp.StatusCode,
		Headers:    resp.Header,
		Body:       respBody,
	}

	switch resp.StatusCode {
	case http.StatusOK, http.StatusCreated, http.StatusAccepted:
		return apiResp, nil
	case http.StatusNoContent:
		return apiResp, nil
	case http.StatusUnauthorized:
		return nil, fmt.Errorf("%w: %s", ErrUnauthorized, string(respBody))
	case http.StatusForbidden:
		return nil, fmt.Errorf("%w: %s", ErrForbidden, string(respBody))
	case http.StatusNotFound:
		return nil, fmt.Errorf("%w: %s", ErrNotFound, string(respBody))
	case http.StatusConflict:
		return nil, fmt.Errorf("%w: %s", ErrConflict, string(respBody))
	case http.StatusUnprocessableEntity:
		return nil, fmt.Errorf("%w: %s", ErrUnprocessableEntity, string(respBody))
	case http.StatusTooManyRequests:
		return nil, fmt.Errorf("%w: %s", ErrRateLimited, string(respBody))
	default:
		return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(respBody))
	}
}

type Response struct {
	StatusCode int
	Headers    http.Header
	Body       []byte
}

func (r *Response) Decode(v interface{}) error {
	if len(r.Body) == 0 {
		return nil
	}
	return json.Unmarshal(r.Body, v)
}

func WithQueryParams(params map[string]string) RequestOption {
	return func(req *http.Request) {
		q := req.URL.Query()
		for k, v := range params {
			q.Set(k, v)
		}
		req.URL.RawQuery = q.Encode()
	}
}

func WithPagination(page, perPage int) RequestOption {
	return func(req *http.Request) {
		q := req.URL.Query()
		if page > 0 {
			q.Set("page", fmt.Sprintf("%d", page))
		}
		if perPage > 0 {
			q.Set("per_page", fmt.Sprintf("%d", perPage))
		}
		req.URL.RawQuery = q.Encode()
	}
}

type PaginatedResponse struct {
	Items      []map[string]interface{}
	Page       int
	PerPage    int
	TotalPages int
	HasNext    bool
}

func ParsePaginationLinks(linkHeader string) (nextPage int, hasNext bool) {
	if linkHeader == "" {
		return 0, false
	}

	links := strings.Split(linkHeader, ",")
	for _, link := range links {
		parts := strings.SplitN(strings.TrimSpace(link), ";", 2)
		if len(parts) != 2 {
			continue
		}

		urlPart := strings.Trim(parts[0], " <>")
		relPart := strings.TrimSpace(parts[1])

		if !strings.Contains(relPart, "rel=\"next\"") {
			continue
		}

		u, err := url.Parse(urlPart)
		if err != nil {
			continue
		}

		pageStr := u.Query().Get("page")
		if pageStr != "" {
			_, _ = fmt.Sscanf(pageStr, "%d", &nextPage)
			hasNext = true
		}
	}

	return nextPage, hasNext
}
