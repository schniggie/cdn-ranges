package provider

import (
	"context"

	cdn_ranges "github.com/taythebot/cdn-ranges"
)

var _ Provider = (*Link11)(nil)

type Link11 struct{}

func (p *Link11) Name() string {
	return "Link11"
}

func (p *Link11) Fetch(ctx context.Context) ([]string, []string, error) {
	return cdn_ranges.ASNPrefixes(ctx, 34309)
}
