"""QuanticData API quickstart — Python (requests).

Scrape one page to Markdown, then run a Google search.
    pip install requests && python3 quickstart.py
"""
import os
import sys

import requests

BASE = "https://api.quanticdata.io/v1"
KEY = os.environ.get("QUANTICDATA_API_KEY")
if not KEY:
    sys.exit("set QUANTICDATA_API_KEY (get one at https://app.quanticdata.io/register)")

session = requests.Session()
session.headers.update({"Authorization": f"Bearer {KEY}", "Content-Type": "application/json"})


def call(path: str, body: dict) -> dict:
    """POST to the API and unwrap the { type, message, payload } envelope."""
    r = session.post(f"{BASE}{path}", json=body, timeout=120)
    data = r.json()
    if data.get("type") == "error" or not r.ok:
        raise RuntimeError(f"{path} failed ({r.status_code}): {data.get('message')}")
    return data.get("payload", {})


def main() -> None:
    page = call("/scrape", {"url": "https://example.com", "format": "markdown"})
    print("--- markdown -------------------------------------------------")
    print(page.get("markdown", "")[:400])
    print("title:", (page.get("metadata") or {}).get("title"))
    print("cost :", (page.get("usage") or {}).get("cost_usd"))

    serp = call("/serp", {"query": "web scraping api", "country": "us", "num": 5})
    print("\n--- organic --------------------------------------------------")
    for row in serp.get("organic", [])[:5]:
        print(f"{row.get('position'):>2}. {row.get('title')}\n    {row.get('link')}")


if __name__ == "__main__":
    main()
