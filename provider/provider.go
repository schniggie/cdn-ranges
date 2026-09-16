package provider

import (
	"context"
)

// Provider interface for all providers
type Provider interface {
	Name() string
	Fetch(ctx context.Context) (ipv4 []string, ipv6 []string, err error)
}

var Providers = []Provider{
	&Akamai{},
	&ArvanCloud{},
	&Azion{},
	&AzureFrontDoor{},
	&Bunny{},
	&CacheFly{},
	&CDN77{},
	&CDNetworks{},
	&CDNvideo{},
	&Cloudflare{},
	&Cloudfront{},
	&DDoSGuard{},
	&Edgecast{},
	&EdgeNext{},
	&Edgio{},
	&F5{},
	&Fastly{},
	&GCore{},
	&GoCache{},
	&Imperva{},
	&Limelight{},
	&Link11{},
	&Medianova{},
	&Myra{},
	&PathNetwork{},
	&Qrator{},
	&Radware{},
	&StormWall{},
	&Sucuri{},
	&X4B{},
}
