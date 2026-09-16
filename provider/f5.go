package provider

import (
	"context"
	"errors"
	"regexp"
	"strings"

	cdn_ranges "github.com/taythebot/cdn-ranges"
)

var _ Provider = (*F5)(nil)

type F5 struct{}

const f5DocURL = "https://docs.cloud.f5.com/docs-v2/platform/reference/network-cloud-ref"

// Sections of the reference page that describe client-facing edge ranges.
// The page also lists log receivers, DNS health checkers and bot defense
// egress addresses, which are not CDN edges.
var f5EdgeSections = []string{
	"Regional Edges",
	"Content Distribution Network",
}

var (
	f5HeadingRe = regexp.MustCompile(`(?is)<h[1-4][^>]*>(.*?)</h[1-4]>`)
	f5CIDRRe    = regexp.MustCompile(`\b(?:\d{1,3}\.){3}\d{1,3}/\d{1,2}\b`)
)

func (p *F5) Name() string {
	return "F5"
}

// Fetch parses the public F5 Distributed Cloud IP reference page. F5 does
// not publish a machine readable list; the page is server rendered so the
// relevant tables can be extracted by heading.
func (p *F5) Fetch(ctx context.Context) ([]string, []string, error) {
	resp, err := cdn_ranges.HttpGet(ctx, f5DocURL)
	if err != nil {
		return nil, nil, err
	}

	page := string(resp)
	headings := f5HeadingRe.FindAllStringSubmatchIndex(page, -1)

	var v4 []string
	for i, h := range headings {
		title := strings.ToLower(stripTags(page[h[2]:h[3]]))
		if !containsAny(title, f5EdgeSections) {
			continue
		}

		end := len(page)
		if i+1 < len(headings) {
			end = headings[i+1][0]
		}

		v4 = append(v4, f5CIDRRe.FindAllString(page[h[1]:end], -1)...)
	}

	if len(v4) == 0 {
		return nil, nil, errors.New("no ranges found on reference page; layout may have changed")
	}

	return v4, nil, nil
}

func stripTags(s string) string {
	return strings.TrimSpace(regexp.MustCompile(`<[^>]+>`).ReplaceAllString(s, ""))
}

func containsAny(s string, needles []string) bool {
	for _, n := range needles {
		if strings.Contains(s, strings.ToLower(n)) {
			return true
		}
	}

	return false
}
