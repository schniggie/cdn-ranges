package provider

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"

	cdn_ranges "github.com/taythebot/cdn-ranges"
)

var _ Provider = (*AzureFrontDoor)(nil)

type AzureFrontDoor struct{}

// Microsoft publishes Azure IP ranges as a weekly "Service Tags" JSON whose
// download URL contains the release date. The stable entry point is the
// download center page, which links to the current file.
const azureServiceTagsPage = "https://www.microsoft.com/en-us/download/confirmation.aspx?id=56519"

var azureServiceTagsURLRe = regexp.MustCompile(`https://download\.microsoft\.com/download/[^"'\s]+/ServiceTags_Public_\d{8}\.json`)

// Frontend is the anycast edge clients connect to, Backend the ranges Front
// Door uses towards origins. Both are Front Door infrastructure.
var azureFrontDoorTags = map[string]bool{
	"AzureFrontDoor.Frontend": true,
	"AzureFrontDoor.Backend":  true,
}

type azureServiceTags struct {
	Values []struct {
		Name       string `json:"name"`
		Properties struct {
			AddressPrefixes []string `json:"addressPrefixes"`
		} `json:"properties"`
	} `json:"values"`
}

func (p *AzureFrontDoor) Name() string {
	return "AzureFrontDoor"
}

func (p *AzureFrontDoor) Fetch(ctx context.Context) ([]string, []string, error) {
	page, err := cdn_ranges.HttpGet(ctx, azureServiceTagsPage)
	if err != nil {
		return nil, nil, err
	}

	url := azureServiceTagsURLRe.FindString(string(page))
	if url == "" {
		return nil, nil, errors.New("could not find ServiceTags_Public download link")
	}

	resp, err := cdn_ranges.HttpGet(ctx, url)
	if err != nil {
		return nil, nil, err
	}

	var tags azureServiceTags
	if err := json.Unmarshal(resp, &tags); err != nil {
		return nil, nil, fmt.Errorf("failed to unmarshal json: %w", err)
	}

	var ranges []string
	found := 0
	for _, v := range tags.Values {
		if azureFrontDoorTags[v.Name] {
			found++
			ranges = append(ranges, v.Properties.AddressPrefixes...)
		}
	}

	if found != len(azureFrontDoorTags) {
		return nil, nil, fmt.Errorf("expected %d AzureFrontDoor service tags, found %d", len(azureFrontDoorTags), found)
	}

	// Sanitize splits the mixed list by address family
	return ranges, nil, nil
}
