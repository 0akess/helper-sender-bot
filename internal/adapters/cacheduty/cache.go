package cacheduty

import (
	"context"
	"regexp"
	"time"

	"log/slog"

	"github.com/maypok86/otter/v2"
)

var (
	dutyRe = regexp.MustCompile(`(?m)^Дежурный:\s+(@\S+)`)
)

type ClientAPI interface {
	ChannelHeader(ctx context.Context, channelID string) (string, error)
}

type Cache struct {
	cache  *otter.Cache[string, string]
	client ClientAPI
	logger *slog.Logger
}

func NewCache(ttl time.Duration, client ClientAPI, logger *slog.Logger) *Cache {
	c := otter.Must[string, string](&otter.Options[string, string]{
		ExpiryCalculator: otter.ExpiryCreating[string, string](ttl),
	})
	return &Cache{
		cache:  c,
		client: client,
		logger: logger,
	}
}

func (c *Cache) GetDutyCache(ctx context.Context, channelID string) (string, error) {
	return c.cache.Get(ctx, channelID,
		otter.LoaderFunc[string, string](func(ctx context.Context, key string) (string, error) {
			header, err := c.client.ChannelHeader(ctx, key)
			if err != nil {
				c.logger.Debug("failed to fetch channel header", "channelID", key, "error", err)
				return "", err
			}
			m := dutyRe.FindStringSubmatch(header)
			if len(m) >= 2 {
				return m[1], nil
			}
			return "", nil
		}),
	)
}
