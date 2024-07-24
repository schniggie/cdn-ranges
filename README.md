# CDN Ranges
Tool to download a list of IP ranges used by CDNs (Cloudflare, Akamai, Incapsula, Fastly, etc). This helps to avoid performing unnecessary port scans when doing bug bounties.

This uses publicly available lists of IP ranges, provided by most providers, and [BGPView](https://bgpview.io/) to query IP ranges for ASNs.

This was heavily inspired by [Project Discovery's cdncheck](https://github.com/projectdiscovery/cdncheck).

## 
Ranges are updated every 60 Minutes and deployed [here](https://github.com/schniggie/cdn-ranges/tree/main/output).

[![Node.js CI](https://github.com/schniggie/cdn-ranges/actions/workflows/node.js.yml/badge.svg)](https://github.com/schniggie/cdn-ranges/actions/workflows/node.js.yml)

## CDN Providers
Provider | ASN or Public List
--- | ---
Apple | "AS714", "AS6185"
Akamai | "AS12222", "AS16625", "AS16702", "AS17204", "AS18680", "AS18717", "AS20189", "AS20940", "AS21342", "AS21357", "AS21399", "AS22207", "AS22452", "AS23454", "AS23455", "AS23903", "AS24319", "AS26008", "AS30675", "AS31107", "AS31108", "AS31109", "AS31110", "AS31377", "AS33047", "AS33905", "AS34164", "AS34850", "AS35204", "AS35993", "AS35994", "AS36183", "AS39836", "AS43639", "AS55409", "AS55770", "AS63949", "AS133103", "AS393560"
CacheFly |  https://cachefly.cachefly.net/ips/rproxy.txt
CDNetworks | "AS36408", "AS38107"
Cloudflare | https://www.cloudflare.com/ips-v4
CloudFront | https://ip-ranges.amazonaws.com/ip-ranges.json
DDoS Guard | "AS57724"
Fastly | https://api.fastly.com/public-ip-list
Incapsula | https://my.incapsula.com/api/integration/v1/ips
Limelight (Edg.io) | "AS22822", "AS23059", "AS26506","AS38622", "AS55429"
MaxCDN | https://support.maxcdn.com/hc/en-us/article_attachments/360051920551/maxcdn_ips.txt
Qrator | "AS200449"
StackPath | "AS12989", "AS18607", "AS20446", "AS33438"
StormWall | "AS59796"
Sucuri | "AS30148"
X4B | "AS136165"


If a provider is missing, please open an issue with a link to their IP ranges or ASN

![](example.png)

## Installation
This code requires Node.js 14.x or higher. I will be providing executables in the near future.

1. Download code
```
git clone https://github.com/taythebot/cdn-ranges
```

2. Install dependencies
```
npm install
```


## Usage
Download ip ranges for all providers
```
node download --output ranges.txt
```

Download for a specific provider (lowercase and replace spaces with -)
```
node download --provider cloudflare
```

Dump in json format
```
node download --format json --output ranges.json
```

Dump in csv format (format: provider,range)
```
node download --format csv --output ranges.csv
```

## Support Formats
* txt (default)
* json
* csv
