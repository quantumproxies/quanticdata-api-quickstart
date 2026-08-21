# QuanticData API quickstart — Ruby (net/http, stdlib only).
#   QUANTICDATA_API_KEY=qd_live_... ruby quickstart.rb
require "json"
require "net/http"
require "uri"

BASE = "https://api.quanticdata.io/v1".freeze

KEY = ENV["QUANTICDATA_API_KEY"]
abort "set QUANTICDATA_API_KEY — https://app.quanticdata.io/register" if KEY.nil? || KEY.empty?

# POST to the API and unwrap the { type, message, payload } envelope.
def qd_call(path, body)
  uri = URI("#{BASE}#{path}")
  http = Net::HTTP.new(uri.host, uri.port)
  http.use_ssl = true
  http.read_timeout = 120

  req = Net::HTTP::Post.new(uri)
  req["Authorization"] = "Bearer #{KEY}"
  req["Content-Type"] = "application/json"
  req.body = JSON.dump(body)

  res = http.request(req)
  data = JSON.parse(res.body) rescue {}
  if res.code.to_i >= 400 || data["type"] == "error"
    raise "#{path} failed (#{res.code}): #{data['message']}"
  end
  data["payload"] || {}
end

page = qd_call("/scrape", { url: "https://example.com", format: "markdown" })
puts "--- markdown -------------------------------------------------"
puts page["markdown"].to_s[0, 400]
puts "title: #{page.dig('metadata', 'title')}"
puts "cost : #{page.dig('usage', 'cost_usd')}"

serp = qd_call("/serp", { query: "web scraping api", country: "us", num: 5 })
puts
puts "--- organic --------------------------------------------------"
serp.fetch("organic", []).first(5).each do |row|
  puts format("%2d. %s\n    %s", row["position"].to_i, row["title"], row["link"])
end
