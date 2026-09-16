package provider

import (
	"context"

	cdn_ranges "github.com/taythebot/cdn-ranges"
)

var _ Provider = (*PathNetwork)(nil)

type PathNetwork struct{}

func (p *PathNetwork) Name() string {
	return "Path-Network"
}

func (p *PathNetwork) Fetch(ctx context.Context) ([]string, []string, error) {
	return cdn_ranges.ASNPrefixes(ctx, 396998)
}
