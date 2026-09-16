package cdn_ranges

import (
	"fmt"
	"net/netip"
	"strings"
)

// Prefixes shorter than these are never legitimate CDN ranges. The
// thresholds leave room for real allocations (Akamai announces a /10 and a
// 2600:1400::/24) while still catching a poisoned or broken source. Tencent
// EdgeOne's deprecated IP list, for example, started returning 0.0.0.0/0.
const (
	MinIPv4PrefixBits = 8
	MinIPv6PrefixBits = 20
)

// Sanitize normalises and validates the ranges returned by a provider.
//
// Entries may be CIDR prefixes or bare addresses (converted to /32 or /128).
// Host bits are masked off, duplicates removed and the result split by
// address family regardless of which list a provider put an entry in. Any
// entry that is not an IP range, or that is implausibly large, is an error
// so a broken source fails the run instead of poisoning the output.
func Sanitize(ranges []string) (v4 []string, v6 []string, err error) {
	seen := make(map[netip.Prefix]struct{}, len(ranges))

	for _, raw := range ranges {
		s := strings.TrimSpace(raw)
		if s == "" {
			continue
		}

		prefix, perr := netip.ParsePrefix(s)
		if perr != nil {
			addr, aerr := netip.ParseAddr(s)
			if aerr != nil {
				return nil, nil, fmt.Errorf("invalid range %q", raw)
			}
			prefix = netip.PrefixFrom(addr, addr.BitLen())
		}
		prefix = prefix.Masked()

		if prefix.Addr().Is4() {
			if prefix.Bits() < MinIPv4PrefixBits {
				return nil, nil, fmt.Errorf("implausibly large range %q", raw)
			}
		} else if prefix.Bits() < MinIPv6PrefixBits {
			return nil, nil, fmt.Errorf("implausibly large range %q", raw)
		}

		if _, dup := seen[prefix]; dup {
			continue
		}
		seen[prefix] = struct{}{}

		if prefix.Addr().Is4() {
			v4 = append(v4, prefix.String())
		} else {
			v6 = append(v6, prefix.String())
		}
	}

	return v4, v6, nil
}
