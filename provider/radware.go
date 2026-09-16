package provider

import (
	"context"

	cdn_ranges "github.com/taythebot/cdn-ranges"
)

var _ Provider = (*Radware)(nil)

type Radware struct{}

func (p *Radware) Name() string {
	return "Radware"
}

func (p *Radware) Fetch(ctx context.Context) ([]string, []string, error) {
	return cdn_ranges.ASNPrefixes(ctx, 25773)
}
