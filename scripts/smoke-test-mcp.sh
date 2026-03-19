#!/usr/bin/env bash

set -euo pipefail

if [[ $# -lt 1 || -z "${1:-}" ]]; then
  echo "Usage: $0 <project-id> [page-size]" >&2
  exit 1
fi

project_id="$1"
page_size="${2:-1}"

if ! [[ "$page_size" =~ ^[0-9]+$ ]] || (( page_size < 1 )); then
  echo "page-size must be a positive integer" >&2
  exit 1
fi

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
binary_path="$(mktemp "${TMPDIR:-/tmp}/mcp-gcp-asset-inventory-smoke.XXXXXX")"

cleanup() {
  rm -f "$binary_path"
}

trap cleanup EXIT

(
  cd "$repo_root"
  go build -o "$binary_path" .
)

response="$(
  printf '%s\n%s\n%s\n' \
    '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2025-11-25","capabilities":{},"clientInfo":{"name":"mcp-smoke-test","version":"0.0.1"}}}' \
    '{"jsonrpc":"2.0","method":"notifications/initialized"}' \
    "{\"jsonrpc\":\"2.0\",\"id\":2,\"method\":\"tools/call\",\"params\":{\"name\":\"list_gcp_assets\",\"arguments\":{\"project_id\":\"${project_id}\",\"page_size\":${page_size}}}}" \
    | "$binary_path"
)"

printf '%s\n' "$response"

if [[ "$response" == *'"error"'* ]]; then
  exit 1
fi
