package provider

import (
	"context"

	cdn_ranges "github.com/taythebot/cdn-ranges"
)

var _ Provider = (*Limelight)(nil)

type Limelight struct{}

// Limelight (later Edgio) wound down its network in 2024/2025, so these ASNs
// may announce no prefixes. They are kept in case routes get re-announced.
var LimelightASNs = []int{
	22822, // LLNW
	23059, // LLNW-LATAM
	26506, // LLNW-SPS
	38622, // LLNW-AU
	55429, // LLNW-IN
}

func (p *Limelight) Name() string {
	return "Limelight"
}

func (p *Limelight) Fetch(ctx context.Context) ([]string, []string, error) {
	return cdn_ranges.ASNsPrefixes(ctx, LimelightASNs...)
}
