package mattermost

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"helper-sender-bot/internal/entity"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	URL, token string
	hc         *http.Client
	log        *slog.Logger
}

func New(url, token string, log *slog.Logger) *Client {
	return &Client{
		URL:   url,
		token: token,
		hc:    &http.Client{Timeout: 10 * time.Second},
		log:   log,
	}
}

func (c *Client) callMM(req *http.Request, out any) (int, error) {
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Accept", "application/json")

	resp, err := c.hc.Do(req)
	if err != nil {
		return 500, fmt.Errorf("http request: %w", err)
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			c.log.Error("failed to close response body", "error", err)
		}
	}()
	c.log.Info("got response", "method", req.Method, "status", resp.Status, "uri", req.URL.RequestURI())

	if resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		c.log.Error(
			"MM API error",
			"method", req.Method,
			"uri", req.URL.RequestURI(),
			"status", resp.StatusCode,
			"body", string(body),
		)
		return resp.StatusCode, fmt.Errorf("%s %s → %s: %s",
			req.Method, req.URL.RequestURI(), resp.Status, strings.TrimSpace(string(body)),
		)
	}

	return resp.StatusCode, json.NewDecoder(resp.Body).Decode(out)
}

func (c *Client) ChannelHeader(ctx context.Context, id string) (string, error) {
	var v struct{ Header string }
	req, _ := http.NewRequestWithContext(ctx, "GET", c.URL+"/api/v4/channels/"+id, nil)
	_, err := c.callMM(req, &v)
	if err != nil {
		return "", err
	}
	return v.Header, nil
}

func (c *Client) fetchPosts(ctx context.Context, path string) ([]entity.Post, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.URL+path, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	var data struct {
		Order []string               `json:"order"`
		Posts map[string]entity.Post `json:"posts"`
	}

	_, err = c.callMM(req, &data)
	if err != nil {
		return nil, err
	}

	out := make([]entity.Post, 0, len(data.Order))
	for _, id := range data.Order {
		out = append(out, data.Posts[id])
	}
	return out, nil
}

func (c *Client) FetchPostsWithSince(ctx context.Context, channelID string, since, perPage int) ([]entity.Post, error) {
	return c.fetchPosts(
		ctx,
		fmt.Sprintf("/api/v4/channels/%s/posts?since=%d&page=0&per_page=%d", channelID, since, perPage),
	)
}

func (c *Client) FetchPostsByPage(ctx context.Context, channelID string, page, perPage int) ([]entity.Post, error) {
	return c.fetchPosts(
		ctx,
		fmt.Sprintf("/api/v4/channels/%s/posts?page=%d&per_page=%d", channelID, page, perPage),
	)
}

func (c *Client) CreatePost(ctx context.Context, channelID, msg, rootID string) (string, int, error) {
	body := struct {
		ChannelID string `json:"channel_id"`
		Message   string `json:"message"`
		RootID    string `json:"root_id,omitempty"`
	}{
		ChannelID: channelID,
		Message:   msg,
		RootID:    rootID,
	}

	payload, err := json.Marshal(body)
	if err != nil {
		return "", 500, fmt.Errorf("marshal body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx,
		http.MethodPost, c.URL+"/api/v4/posts", bytes.NewReader(payload))
	if err != nil {
		return "", 500, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	var resp struct {
		ID string `json:"id"`
	}
	statusCode, err := c.callMM(req, &resp)
	if err != nil {
		return "", statusCode, err
	}
	return resp.ID, statusCode, nil
}
