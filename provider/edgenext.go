package provider

import (
	"context"

	cdn_ranges "github.com/taythebot/cdn-ranges"
)

var _ Provider = (*EdgeNext)(nil)

type EdgeNext struct{}

var EdgeNextASNs = []int{
	139057,
	149981,
}

func (p *EdgeNext) Name() string {
	return "EdgeNext"
}

func (p *EdgeNext) Fetch(ctx context.Context) ([]string, []string, error) {
	return cdn_ranges.ASNsPrefixes(ctx, EdgeNextASNs...)
}
