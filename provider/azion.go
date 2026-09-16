package provider

import (
	"context"

	cdn_ranges "github.com/taythebot/cdn-ranges"
)

var _ Provider = (*Azion)(nil)

type Azion struct{}

func (p *Azion) Name() string {
	return "Azion"
}

func (p *Azion) Fetch(ctx context.Context) ([]string, []string, error) {
	return cdn_ranges.ASNPrefixes(ctx, 52580)
}
