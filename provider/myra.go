package provider

import (
	"context"

	cdn_ranges "github.com/taythebot/cdn-ranges"
)

var _ Provider = (*Myra)(nil)

type Myra struct{}

func (p *Myra) Name() string {
	return "Myra"
}

func (p *Myra) Fetch(ctx context.Context) ([]string, []string, error) {
	return cdn_ranges.ASNPrefixes(ctx, 41179)
}
