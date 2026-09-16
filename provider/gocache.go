package provider

import (
	"context"
	"strings"

	cdn_ranges "github.com/taythebot/cdn-ranges"
)

var _ Provider = (*GoCache)(nil)

type GoCache struct{}

func (p *GoCache) Name() string {
	return "GoCache"
}

func (p *GoCache) Fetch(ctx context.Context) ([]string, []string, error) {
	resp, err := cdn_ranges.HttpGet(ctx, "https://gocache.com.br/ips")
	if err != nil {
		return nil, nil, err
	}

	return strings.Fields(string(resp)), nil, nil
}
