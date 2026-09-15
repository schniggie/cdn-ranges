# CDN Ranges

Tool to download a list of IPv4 and IPv6 ranges used by CDNs. This helps to avoid performing unnecessary port scans when
doing bug bounties.

This uses publicly available lists of IP ranges, provided by most providers, and [IP Guide](https://ip.guide) to query IP ranges for ASNs.

This is a fork of [taythebot/cdn-ranges](https://github.com/taythebot/cdn-ranges) that additionally publishes the
generated lists in this repository.

## Pre-generated lists

[![Update CDN ranges](https://github.com/schniggie/cdn-ranges/actions/workflows/update-ranges.yml/badge.svg)](https://github.com/schniggie/cdn-ranges/actions/workflows/update-ranges.yml)

Ranges are refreshed every 60 minutes and committed to [`output/`](https://github.com/schniggie/cdn-ranges/tree/main/output):

| File                        | Content                                        |
|-----------------------------|------------------------------------------------|
| `output/all-cdn.txt`        | IPv4 ranges of all providers                   |
| `output/all-cdn-ipv6.txt`   | IPv6 ranges of all providers                   |
| `output/all-cdn.csv`        | `provider,type,range` for all providers        |
| `output/<provider>.txt`     | IPv4 ranges of a single provider               |
| `output/<provider>-ipv6.txt`| IPv6 ranges of a single provider               |
| `output/.lastrun`           | UTC timestamp of the last run                  |

Raw file example: `https://raw.githubusercontent.com/schniggie/cdn-ranges/main/output/all-cdn.txt`

## CDN Providers

| Provider   | ASN or Public List                                                                                                                                                                                                                                                                                                                                              |
|------------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| Akamai     | AS12222, AS16625, AS16702, AS17204, AS18680, AS18717, AS20189, AS20940, AS21342, AS21357, AS21399, AS22207, AS22452, AS23454, AS23455, AS23903, AS24319, AS26008, AS30675, AS31107, AS31108, AS31109, AS31110, AS31377, AS33047, AS33905, AS34164, AS34850, AS35204, AS35993, AS35994, AS36183, AS39836, AS43639, AS55409, AS55770, AS63949, AS133103, AS393560 |
| ArvanCloud | https://www.arvancloud.ir/en/ips.txt                                                                                                                                                                                                                                                                                                                            |
| Bunny      | AS200325                                                                                                                                                                                                                                                                                                                                                        |
| CacheFly   | https://cachefly.cachefly.net/ips/rproxy.txt                                                                                                                                                                                                                                                                                                                    |
| CDN77      | https://prefixlists.tools.cdn77.com/public_lmax_prefixes.json                                                                                                                                                                                                                                                                                                   |
| CDNetworks | AS36408, AS38107                                                                                                                                                                                                                                                                                                                                                |
| Cloudflare | https://www.cloudflare.com/ips-v4 https://www.cloudflare.com/ips-v6                                                                                                                                                                                                                                                                                             |
| CloudFront | https://ip-ranges.amazonaws.com/ip-ranges.json                                                                                                                                                                                                                                                                                                                  |
| DDoS-Guard | AS57724                                                                                                                                                                                                                                                                                                                                                         |
| Edgecast   | AS15133                                                                                                                                                                                                                                                                                                                                                         |
| EdgeNext   | AS139057, AS149981                                                                                                                                                                                                                                                                                                                                              |
| Edgio      | AS60261                                                                                                                                                                                                                                                                                                                                                         |
| F5         | https://docs.cloud.f5.com/docs-v2/platform/reference/network-cloud-ref                                                                                                                                                                                                                                                                                          |
| Fastly     | https://api.fastly.com/public-ip-list                                                                                                                                                                                                                                                                                                                           |
| Gcore      | https://api.gcore.com/cdn/public-ip-list                                                                                                                                                                                                                                                                                                                        |
| Imperva    | https://my.imperva.com/api/integration/v1/ips                                                                                                                                                                                                                                                                                                                   |
| Leaseweb   | https://networksdb.io/ip-addresses-of/leaseweb-cdn-bv                                                                                                                                                                                                                                                                                                           |
| Limelight  | AS22822, AS23059, AS26506, AS38622, AS55429                                                                                                                                                                                                                                                                                                                     |
| Medianova  | https://cloud.medianova.com/api/v1/ip/blocks-list                                                                                                                                                                                                                                                                                                               |
| Qrator     | AS200449, AS209671                                                                                                                                                                                                                                                                                                                                              |
| StormWall  | AS59796                                                                                                                                                                                                                                                                                                                                                         |
| Sucuri     | AS30148                                                                                                                                                                                                                                                                                                                                                         |
| X4B        | AS136165                                                                                                                                                                                                                                                                                                                                                        |

Note: Edgecast, Edgio and Limelight (all part of the former Edgio network) currently announce no prefixes and produce
empty lists. StackPath and MaxCDN were removed because they shut down; their former ASNs now belong to unrelated
organisations.

If a provider is missing, please open an issue with a link to their IP ranges or ASN.

## Installation

Requires Go 1.23 or newer.

```bash
git clone https://github.com/schniggie/cdn-ranges
cd cdn-ranges
go build -o cdn-ranges ./cmd/cdn-ranges
```

## Usage

Download IPv4 and IPv6 ranges for all providers

```bash
cdn-ranges -output ranges.txt
```

Download IPv4 ranges only

```bash
cdn-ranges -ipv4 -output ranges.txt
```

Download IPv6 ranges only

```bash
cdn-ranges -ipv6 -output ranges.txt
```

Download for a specific provider

```bash
cdn-ranges -provider cloudflare
```

List available providers

```bash
cdn-ranges -list
```

Dump in json format

```bash
cdn-ranges -format json -output ranges.json
```

Dump in csv format (provider,type,range)

```bash
cdn-ranges -format csv -output ranges.csv
```

Regenerate the pre-generated lists in `output/` locally

```bash
scripts/update-output.sh
```

## Output Formats

* txt (default)
* json
* csv

## Adding a provider

1. Create `provider/<name>.go` implementing the `Provider` interface (`Name()` and `Fetch(ctx)`).
   * For ASN based providers return `cdn_ranges.ASNsPrefixes(ctx, <asn>, ...)`.
   * For published lists fetch them with `cdn_ranges.HttpGet(ctx, url)`, which validates the HTTP status and retries
     transient failures.
2. Register the provider in the `Providers` slice in `provider/provider.go`.
3. Add it to the provider table above.

The CI workflow runs `gofmt`, `go vet` and a build on every pull request; the hourly update workflow picks up the new
provider automatically via `cdn-ranges -list`.
