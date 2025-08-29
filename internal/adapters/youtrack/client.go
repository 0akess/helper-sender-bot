package youtrack

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"log/slog"
)

type Issue struct {
	ID      string    `json:"idReadable"`
	Summary string    `json:"summary"`
	Created time.Time `json:"created"`
	State   string
}

type Client struct {
	Host       string
	Token      string
	httpClient *http.Client
	log        *slog.Logger
}

func NewYouTrackClient(host, token string, timeout time.Duration, log *slog.Logger) *Client {
	return &Client{
		Host:       host,
		Token:      token,
		httpClient: &http.Client{Timeout: timeout},
		log:        log,
	}
}

func (c *Client) FetchOpenBugs(project string) ([]Issue, error) {
	jql := fmt.Sprintf("project: %s AND Type: Bug AND State: Open", project)
	return c.FetchIssues(jql, []string{"idReadable", "summary", "created", "state(name)"})
}

func (c *Client) FetchIssues(jql string, fields []string) ([]Issue, error) {
	params := url.Values{}
	params.Set("query", jql)
	params.Set("fields", strings.Join(fields, ","))

	var raw []struct {
		ID      string                `json:"idReadable"`
		Summary string                `json:"summary"`
		Created time.Time             `json:"created"`
		State   struct{ Name string } `json:"state"`
	}
	if err := c.getJSON("/api/issues", params, &raw); err != nil {
		return nil, err
	}

	issues := make([]Issue, len(raw))
	for i, r := range raw {
		issues[i] = Issue{
			ID:      r.ID,
			Summary: r.Summary,
			Created: r.Created,
			State:   r.State.Name,
		}
	}
	return issues, nil
}

func (c *Client) GetIssueState(issueID string) (string, error) {
	params := url.Values{}
	params.Set("fields", "state(name)")

	var wrap struct {
		State struct{ Name string } `json:"state"`
	}
	if err := c.getJSON(fmt.Sprintf("/api/issues/%s", issueID), params, &wrap); err != nil {
		return "", err
	}
	return wrap.State.Name, nil
}

func (c *Client) getJSON(path string, params url.Values, v interface{}) error {
	endpoint := fmt.Sprintf("%s%s", c.Host, path)
	req, err := http.NewRequest("GET", endpoint+"?"+params.Encode(), nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+c.Token)
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			c.log.Error("failed to close response body", "error", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status: %s", resp.Status)
	}
	return json.NewDecoder(resp.Body).Decode(v)
}
