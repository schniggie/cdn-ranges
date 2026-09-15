#!/usr/bin/env bash
# Regenerates the published range lists in output/.
#
# Layout (kept backwards compatible with the original Node.js based lists):
#   output/all-cdn.txt          IPv4 ranges of all providers
#   output/all-cdn-ipv6.txt     IPv6 ranges of all providers
#   output/all-cdn.csv          provider,type,range for all providers
#   output/<provider>.txt       IPv4 ranges of a single provider
#   output/<provider>-ipv6.txt  IPv6 ranges of a single provider
#   output/.lastrun             timestamp of the last successful run
#
# Usage: scripts/update-output.sh [path-to-cdn-ranges-binary]
set -euo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
bin="${1:-}"
out="${repo_root}/output"

if [[ -z "${bin}" ]]; then
  bin="$(mktemp -d)/cdn-ranges"
  (cd "${repo_root}" && go build -o "${bin}" ./cmd/cdn-ranges)
fi

mkdir -p "${out}"

# Fetch everything exactly once; all other files are derived from this
# snapshot so they are guaranteed to be consistent with each other.
"${bin}" -format csv -output "${out}/all-cdn.csv"

awk -F, 'NR > 1 && $2 == "ipv4" { print $3 }' "${out}/all-cdn.csv" > "${out}/all-cdn.txt"
awk -F, 'NR > 1 && $2 == "ipv6" { print $3 }' "${out}/all-cdn.csv" > "${out}/all-cdn-ipv6.txt"

while read -r provider; do
  awk -F, -v p="${provider}" 'NR > 1 && tolower($1) == p && $2 == "ipv4" { print $3 }' \
    "${out}/all-cdn.csv" > "${out}/${provider}.txt"
  awk -F, -v p="${provider}" 'NR > 1 && tolower($1) == p && $2 == "ipv6" { print $3 }' \
    "${out}/all-cdn.csv" > "${out}/${provider}-ipv6.txt"
done < <("${bin}" -list)

date -u +"%Y-%m-%dT%H:%M:%SZ" > "${out}/.lastrun"
