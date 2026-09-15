package provider

import (
	"context"
	"strings"

	cdn_ranges "github.com/taythebot/cdn-ranges"
)

var _ Provider = (*ArvanCloud)(nil)

type ArvanCloud struct{}

func (p *ArvanCloud) Name() string {
	return "ArvanCloud"
}

func (p *ArvanCloud) Fetch(ctx context.Context) ([]string, []string, error) {
	resp, err := cdn_ranges.HttpGet(ctx, "https://www.arvancloud.ir/en/ips.txt")
	if err != nil {
		return nil, nil, err
	}

	return strings.Fields(string(resp)), nil, nil
}
