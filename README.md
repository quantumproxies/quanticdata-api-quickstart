# QuanticData API quickstart — one Bearer key, eight languages

Minimal, copy-paste first calls to the [QuanticData](https://quanticdata.io) Data APIs in
**Python, Node.js, Go, PHP, Ruby, Java, C# and plain curl**. Every file does the same thing —
scrape a page to Markdown, then run a Google search — so you can diff the language you know
against the one you need.

- Base URL: `https://api.quanticdata.io/v1`
- Auth: `Authorization: Bearer $QUANTICDATA_API_KEY`
- One envelope for every endpoint: `{ "type": "response", "message": "...", "payload": { ... } }`
- Failures return `{ "type": "error", "message": "..." }` and cost nothing (pay per success)

Full reference: **[quanticdata.io/docs](https://quanticdata.io/docs/)**.

## The six endpoints

| Endpoint | What it does | Price | Docs |
|---|---|---|---|
| `POST /v1/scrape` | one page → Markdown, HTML or text | $0.0002 | [Web Scraping API](https://quanticdata.io/web-scraping-api/) |
| `POST /v1/serp` | search results as structured JSON | from $0.0005 | [SERP API](https://quanticdata.io/serp-api/) |
| `POST /v1/map` | every URL of a site, from sitemaps + links | $0.0005 | [Crawl & Map](https://quanticdata.io/crawl-map/) |
| `POST /v1/crawl` | async BFS crawl; poll `GET /v1/crawl/:jobId` | $0.0003/page | [Crawl & Map](https://quanticdata.io/crawl-map/) |
| `POST /v1/batch` | up to 1,000 known URLs; poll `GET /v1/batch/:jobId` | $0.0002/URL | [Web Scraping API](https://quanticdata.io/web-scraping-api/) |
| `POST /v1/seo-audit` | no-JS vs rendered view of a URL, diffed | $0.0012 | [SEO Audit](https://quanticdata.io/seo-audit/) |

Ready-made scrapers with semantic inputs (keyword + city, ASIN, domain…) live behind
`/v1/scraper/collectors` — see [Collectors](https://quanticdata.io/collectors/).

## Get a key

Sign up at [app.quanticdata.io](https://app.quanticdata.io/register) — keys look like
`qd_live_…`, the pay-as-you-go tier includes a **free monthly allowance** with no card, and the
rate limit is 60 requests/minute. The current allowance and unit prices are on
[quanticdata.io/docs](https://quanticdata.io/docs/).

```bash
export QUANTICDATA_API_KEY=qd_live_your_key_here
```

## Run it

```bash
python3 quickstart.py                  # requests
node quickstart.mjs                     # no dependencies, native fetch
go run quickstart.go                   # stdlib only
php quickstart.php                     # curl extension
ruby quickstart.rb                     # net/http
java Quickstart.java                   # JDK 11+ single-file mode
dotnet run                             # see Quickstart.cs
bash quickstart.sh                     # curl + jq
```

## The envelope, once

Everything under `payload`. For `/v1/scrape` the useful keys are:

```jsonc
{
  "type": "response",
  "message": "Scrape complete",
  "payload": {
    "url": "https://example.com",
    "markdown": "# Example Domain\n…",
    "metadata": { "title": "Example Domain", "description": null, "canonical": "…" },
    "usage": { "cost_usd": 0.0002 }
  }
}
```

So the one line you need in any language is: *take `payload`, ignore the rest.*

## Errors worth handling

| HTTP | Meaning | What to do |
|---|---|---|
| 401 | missing/invalid `Authorization` | check the `Bearer ` prefix |
| 402 | out of balance or over budget | top up in the dashboard |
| 429 | over 60 req/min | back off, then retry |
| 5xx | target fought back after retries | the call is free; retry later |

## More

- [Documentation](https://quanticdata.io/docs/) · [Collectors catalogue](https://quanticdata.io/collectors/) · [MCP server](https://quanticdata.io/mcp-server/)
- [Web Data API for AI](https://quanticdata.io/web-data-api-for-ai/) · [AI web scraping](https://quanticdata.io/ai-web-scraping-service/)
- Free tools: [curl converter](https://quanticdata.io/tools/curl-converter/) · [robots.txt tester](https://quanticdata.io/tools/robots-txt-tester/)

MIT licensed.
