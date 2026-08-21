/**
 * QuanticData API quickstart — Node.js 18+ (native fetch, zero dependencies).
 *   QUANTICDATA_API_KEY=qd_live_... node quickstart.mjs
 */
const BASE = "https://api.quanticdata.io/v1";
const KEY = process.env.QUANTICDATA_API_KEY;
if (!KEY) {
  console.error("set QUANTICDATA_API_KEY (get one at https://app.quanticdata.io/register)");
  process.exit(1);
}

/** POST to the API and unwrap the { type, message, payload } envelope. */
async function call(path, body) {
  const res = await fetch(BASE + path, {
    method: "POST",
    headers: { Authorization: `Bearer ${KEY}`, "Content-Type": "application/json" },
    body: JSON.stringify(body),
  });
  const data = await res.json();
  if (!res.ok || data.type === "error") {
    throw new Error(`${path} failed (${res.status}): ${data.message}`);
  }
  return data.payload ?? {};
}

const page = await call("/scrape", { url: "https://example.com", format: "markdown" });
console.log("--- markdown -------------------------------------------------");
console.log((page.markdown ?? "").slice(0, 400));
console.log("title:", page.metadata?.title);
console.log("cost :", page.usage?.cost_usd);

const serp = await call("/serp", { query: "web scraping api", country: "us", num: 5 });
console.log("\n--- organic --------------------------------------------------");
for (const row of (serp.organic ?? []).slice(0, 5)) {
  console.log(`${String(row.position).padStart(2)}. ${row.title}\n    ${row.link}`);
}
