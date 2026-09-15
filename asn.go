package cdn_ranges

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"go.uber.org/ratelimit"
)

var (
	rateLimiterOnce sync.Once
	rateLimiter     ratelimit.Limiter
)

type IPGuideResponse struct {
	Asn          int    `json:"asn"`
	Name         string `json:"name"`
	Organization string `json:"organization"`
	Country      string `json:"country"`
	Rir          string `json:"rir"`
	Routes       struct {
		V4 []string `json:"v4"`
		V6 []string `json:"v6"`
	} `json:"routes"`
}

// GetRateLimit returns a global reusable rate limiter
func GetRateLimit() ratelimit.Limiter {
	rateLimiterOnce.Do(func() {
		rateLimiter = ratelimit.New(3)
	})

	return rateLimiter
}

// ASNPrefixes fetches IPv4 and IPv6 prefixes for a given ASN
func ASNPrefixes(ctx context.Context, asn int) ([]string, []string, error) {
	GetRateLimit().Take()

	resp, err := HttpGet(ctx, fmt.Sprintf("https://ip.guide/AS%d", asn))
	if err != nil {
		return nil, nil, err
	}

	var response IPGuideResponse
	if err := json.Unmarshal(resp, &response); err != nil {
		return nil, nil, fmt.Errorf("failed to perform unmarshal json: %w", err)
	}

	return response.Routes.V4, response.Routes.V6, nil
}

// ASNsPrefixes fetches and concatenates the IPv4 and IPv6 prefixes of multiple ASNs
func ASNsPrefixes(ctx context.Context, asns ...int) ([]string, []string, error) {
	var (
		v4 []string
		v6 []string
	)
	for _, asn := range asns {
		asnV4, asnV6, err := ASNPrefixes(ctx, asn)
		if err != nil {
			return nil, nil, fmt.Errorf("AS%d: %w", asn, err)
		}

		v4 = append(v4, asnV4...)
		v6 = append(v6, asnV6...)
	}

	return v4, v6, nil
}
