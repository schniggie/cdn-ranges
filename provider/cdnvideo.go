package provider

import (
	"context"

	cdn_ranges "github.com/taythebot/cdn-ranges"
)

var _ Provider = (*CDNvideo)(nil)

type CDNvideo struct{}

func (p *CDNvideo) Name() string {
	return "CDNvideo"
}

func (p *CDNvideo) Fetch(ctx context.Context) ([]string, []string, error) {
	return cdn_ranges.ASNPrefixes(ctx, 57363)
}
