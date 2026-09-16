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

# Fetch everything exactly once into a temporary file; all other files are
# derived from this snapshot so they are guaranteed to be consistent with
# each other. The published files are only replaced once the snapshot has
# passed the regression checks below.
new_csv="$(mktemp)"
trap 'rm -f "${new_csv}"' EXIT
"${bin}" -format csv -output "${new_csv}"

count_provider() { # <csv> <provider, any case>
  awk -F, -v p="$2" 'NR > 1 && tolower($1) == tolower(p) { n++ } END { print n + 0 }' "$1"
}
count_all() {
  awk 'NR > 1 { n++ } END { print n + 0 }' "$1"
}

# Regression guard: a provider that used to have ranges but now has none,
# or a total that shrank by more than 30 %, points at a broken source rather
# than a real change. Abort instead of publishing a degraded list.
if [[ -s "${out}/all-cdn.csv" ]]; then
  old_total="$(count_all "${out}/all-cdn.csv")"
  new_total="$(count_all "${new_csv}")"
  if (( new_total * 10 < old_total * 7 )); then
    echo "ERROR: total ranges dropped from ${old_total} to ${new_total}" >&2
    exit 1
  fi

  while read -r provider; do
    old_n="$(count_provider "${out}/all-cdn.csv" "${provider}")"
    new_n="$(count_provider "${new_csv}" "${provider}")"
    if (( old_n > 0 && new_n == 0 )); then
      echo "ERROR: ${provider} previously had ${old_n} ranges, now none" >&2
      exit 1
    fi
  done < <("${bin}" -list)
fi

mv "${new_csv}" "${out}/all-cdn.csv"
trap - EXIT

awk -F, 'NR > 1 && $2 == "ipv4" { print $3 }' "${out}/all-cdn.csv" > "${out}/all-cdn.txt"
awk -F, 'NR > 1 && $2 == "ipv6" { print $3 }' "${out}/all-cdn.csv" > "${out}/all-cdn-ipv6.txt"

while read -r provider; do
  file="${out}/${provider,,}"
  awk -F, -v p="${provider}" 'NR > 1 && tolower($1) == tolower(p) && $2 == "ipv4" { print $3 }' \
    "${out}/all-cdn.csv" > "${file}.txt"
  awk -F, -v p="${provider}" 'NR > 1 && tolower($1) == tolower(p) && $2 == "ipv6" { print $3 }' \
    "${out}/all-cdn.csv" > "${file}-ipv6.txt"
done < <("${bin}" -list)

date -u +"%Y-%m-%dT%H:%M:%SZ" > "${out}/.lastrun"
