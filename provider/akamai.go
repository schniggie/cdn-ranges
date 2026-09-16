package provider

import (
	"context"

	cdn_ranges "github.com/taythebot/cdn-ranges"
)

var _ Provider = (*Akamai)(nil)

type Akamai struct{}

var AkamaiASNs = []int{
	12222,
	16625,
	16702,
	17204,
	18680,
	18717,
	20189,
	20940,
	21342,
	21357,
	21399,
	22207,
	22452,
	23454,
	23455,
	23903,
	24319,
	26008,
	30675,
	31107,
	31108,
	31109,
	31110,
	31377,
	33047,
	33905,
	34164,
	34850,
	35204,
	35993,
	35994,
	36183,
	39836,
	43639,
	55409,
	55770,
	63949,
	133103,
	393560,
	// Prolexic DDoS mitigation network (Akamai)
	32787,
	49846,
	213120,
}

func (p *Akamai) Name() string {
	return "Akamai"
}

func (p *Akamai) Fetch(ctx context.Context) ([]string, []string, error) {
	return cdn_ranges.ASNsPrefixes(ctx, AkamaiASNs...)
}
