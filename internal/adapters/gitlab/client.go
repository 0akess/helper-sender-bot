package gitlab

import (
	"context"
	"fmt"
	"gitlab.com/gitlab-org/api/client-go"
	"helper-sender-bot/internal/entity"
	"strings"
)

var clients = map[string]*Client{}

type Client struct {
	gl *gitlab.Client
}

type GitConfigs struct {
	BaseURL string
	Token   string
}

func New(gits []GitConfigs) (*Client, error) {
	for _, cfg := range gits {
		gl, err := gitlab.NewClient(cfg.Token, gitlab.WithBaseURL(cfg.BaseURL))
		if err != nil {
			return nil, fmt.Errorf("gitlab client: %w", err)
		}
		clients[cfg.BaseURL] = &Client{gl: gl}
	}
	return &Client{}, nil
}

func (c *Client) GetMRInfo(
	ctx context.Context,
	p entity.MergeRequestPayload,
	baseURL string,
) (entity.MergeRequestInfo, error) {
	client, ok := clients[baseURL]
	if !ok {
		return entity.MergeRequestInfo{}, fmt.Errorf("client not found for %s", baseURL)
	}
	diffs, _, err := client.gl.MergeRequests.ListMergeRequestDiffs(p.ProjectID, p.MRIID, nil, gitlab.WithContext(ctx))
	if err != nil {
		return entity.MergeRequestInfo{}, fmt.Errorf("fetch MR diffs: %w", err)
	}

	total, hasTests := getMetaInfo(diffs)

	return entity.MergeRequestInfo{TotalLinesChanged: total, HasTestChanges: hasTests}, nil
}

func getMetaInfo(diffs []*gitlab.MergeRequestDiff) (totalChange int, hasTests bool) {
	for _, d := range diffs {
		path := strings.ToLower(d.NewPath)
		if strings.Contains(path, "test") {
			hasTests = true
		}
		add, del := diffStats(d.Diff)
		totalChange += add + del
	}
	return totalChange, hasTests
}

func diffStats(diff string) (add, del int) {
	for _, l := range strings.Split(diff, "\n") {
		switch {
		case strings.HasPrefix(l, "+") && !strings.HasPrefix(l, "+++"):
			add++
		case strings.HasPrefix(l, "-") && !strings.HasPrefix(l, "---"):
			del++
		}
	}
	return
}
