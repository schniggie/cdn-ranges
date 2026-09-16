package provider

import (
	"context"

	cdn_ranges "github.com/taythebot/cdn-ranges"
)

var _ Provider = (*Qrator)(nil)

type Qrator struct{}

var QratorASNs = []int{
	200449, // QRATOR-CZ
	209671, // QRATOR-SW
	211112, // QRATOR-FZ
	197068, // HLL (Qrator)
	51115,  // HLL-AS
}

func (p *Qrator) Name() string {
	return "Qrator"
}

func (p *Qrator) Fetch(ctx context.Context) ([]string, []string, error) {
	return cdn_ranges.ASNsPrefixes(ctx, QratorASNs...)
}
