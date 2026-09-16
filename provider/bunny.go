package provider

import (
	"context"
	"encoding/json"
	"fmt"

	cdn_ranges "github.com/taythebot/cdn-ranges"
)

var _ Provider = (*Bunny)(nil)

type Bunny struct{}

func (p *Bunny) Name() string {
	return "Bunny"
}

// Bunny runs most of its edge on partner networks, so its own ASN covers
// only a small part of the edge. The official edge server list is the
// authoritative source; the ASN is added for Bunny's own address space.
func (p *Bunny) Fetch(ctx context.Context) ([]string, []string, error) {
	v4, err := bunnyEdgeList(ctx, "https://api.bunny.net/system/edgeserverlist")
	if err != nil {
		return nil, nil, err
	}

	v6, err := bunnyEdgeList(ctx, "https://api.bunny.net/system/edgeserverlist/ipv6")
	if err != nil {
		return nil, nil, err
	}

	asnV4, asnV6, err := cdn_ranges.ASNPrefixes(ctx, 200325)
	if err != nil {
		return nil, nil, err
	}

	return append(v4, asnV4...), append(v6, asnV6...), nil
}

func bunnyEdgeList(ctx context.Context, url string) ([]string, error) {
	resp, err := cdn_ranges.HttpGet(ctx, url)
	if err != nil {
		return nil, err
	}

	var ips []string
	if err := json.Unmarshal(resp, &ips); err != nil {
		return nil, fmt.Errorf("failed to unmarshal json: %w", err)
	}

	return ips, nil
}
