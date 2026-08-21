// QuanticData API quickstart — C# / .NET 8 (System.Net.Http + System.Text.Json).
//   QUANTICDATA_API_KEY=qd_live_... dotnet run
using System.Text;
using System.Text.Json;

const string Base = "https://api.quanticdata.io/v1";

var key = Environment.GetEnvironmentVariable("QUANTICDATA_API_KEY");
if (string.IsNullOrWhiteSpace(key))
{
    Console.Error.WriteLine("set QUANTICDATA_API_KEY — https://app.quanticdata.io/register");
    return 1;
}

using var http = new HttpClient { Timeout = TimeSpan.FromSeconds(120) };
http.DefaultRequestHeaders.Add("Authorization", $"Bearer {key}");

// POST to the API and unwrap the { type, message, payload } envelope.
async Task<JsonElement> CallAsync(string path, object body)
{
    var json = JsonSerializer.Serialize(body);
    using var content = new StringContent(json, Encoding.UTF8, "application/json");
    using var res = await http.PostAsync(Base + path, content);
    var text = await res.Content.ReadAsStringAsync();

    using var doc = JsonDocument.Parse(text);
    var root = doc.RootElement;
    var type = root.TryGetProperty("type", out var t) ? t.GetString() : null;
    if (!res.IsSuccessStatusCode || type == "error")
    {
        var msg = root.TryGetProperty("message", out var m) ? m.GetString() : text;
        throw new InvalidOperationException($"{path} failed ({(int)res.StatusCode}): {msg}");
    }
    return root.GetProperty("payload").Clone();
}

var page = await CallAsync("/scrape", new { url = "https://example.com", format = "markdown" });
var markdown = page.TryGetProperty("markdown", out var md) ? md.GetString() ?? "" : "";
Console.WriteLine("--- markdown -------------------------------------------------");
Console.WriteLine(markdown[..Math.Min(400, markdown.Length)]);

var serp = await CallAsync("/serp", new { query = "web scraping api", country = "us", num = 5 });
Console.WriteLine();
Console.WriteLine("--- organic --------------------------------------------------");
if (serp.TryGetProperty("organic", out var organic))
{
    foreach (var row in organic.EnumerateArray())
    {
        var position = row.TryGetProperty("position", out var p) ? p.ToString() : "?";
        var title = row.TryGetProperty("title", out var ti) ? ti.GetString() : "";
        var link = row.TryGetProperty("link", out var l) ? l.GetString() : "";
        Console.WriteLine($"{position,2}. {title}\n    {link}");
    }
}

return 0;
