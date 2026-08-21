#!/usr/bin/env bash
# QuanticData API quickstart — curl + jq.
#   QUANTICDATA_API_KEY=qd_live_... bash quickstart.sh
set -euo pipefail
: "${QUANTICDATA_API_KEY:?set QUANTICDATA_API_KEY — https://app.quanticdata.io/register}"

BASE=https://api.quanticdata.io/v1
AUTH=(-H "Authorization: Bearer $QUANTICDATA_API_KEY" -H "Content-Type: application/json")

echo "--- markdown -------------------------------------------------"
curl -sS "$BASE/scrape" "${AUTH[@]}" \
  -d '{"url":"https://example.com","format":"markdown"}' \
| jq -r '.payload.markdown' | head -20

echo
echo "--- organic --------------------------------------------------"
curl -sS "$BASE/serp" "${AUTH[@]}" \
  -d '{"query":"web scraping api","country":"us","num":5}' \
| jq -r '.payload.organic[] | "\(.position). \(.title)\n   \(.link)"'
