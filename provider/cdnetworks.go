package provider

import (
	"context"

	cdn_ranges "github.com/taythebot/cdn-ranges"
)

var _ Provider = (*CDNetworks)(nil)

type CDNetworks struct{}

var CDNetworksASNs = []int{
	36408,  // CDNETWORKSUS
	38107,  // CDNETWORKS-AS-KR
	204720, // CDNetworks Global Cloud Network
}

func (p *CDNetworks) Name() string {
	return "CDNetworks"
}

func (p *CDNetworks) Fetch(ctx context.Context) ([]string, []string, error) {
	return cdn_ranges.ASNsPrefixes(ctx, CDNetworksASNs...)
}
