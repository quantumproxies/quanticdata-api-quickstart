// QuanticData API quickstart — Java 11+ (single-file source mode, no build tool).
//   QUANTICDATA_API_KEY=qd_live_... java Quickstart.java
import java.net.URI;
import java.net.http.HttpClient;
import java.net.http.HttpRequest;
import java.net.http.HttpResponse;
import java.time.Duration;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

public class Quickstart {

    private static final String BASE = "https://api.quanticdata.io/v1";
    private static final HttpClient CLIENT =
            HttpClient.newBuilder().connectTimeout(Duration.ofSeconds(20)).build();

    /** POST a JSON body and return the raw response text. */
    private static String call(String path, String json) throws Exception {
        String key = System.getenv("QUANTICDATA_API_KEY");
        HttpRequest req = HttpRequest.newBuilder(URI.create(BASE + path))
                .header("Authorization", "Bearer " + key)
                .header("Content-Type", "application/json")
                .timeout(Duration.ofSeconds(120))
                .POST(HttpRequest.BodyPublishers.ofString(json))
                .build();

        HttpResponse<String> res = CLIENT.send(req, HttpResponse.BodyHandlers.ofString());
        if (res.statusCode() >= 400) {
            throw new IllegalStateException(path + " failed (" + res.statusCode() + "): " + res.body());
        }
        return res.body();
    }

    /**
     * The JDK ships no JSON parser, and this file deliberately has no dependencies —
     * so we pull the two fields we print with a regex. In real code use Jackson or Gson
     * and map the { type, message, payload } envelope to a record.
     */
    private static String field(String json, String name) {
        Matcher m = Pattern.compile("\"" + name + "\"\\s*:\\s*\"((?:[^\"\\\\]|\\\\.)*)\"").matcher(json);
        return m.find() ? m.group(1).replace("\\n", "\n").replace("\\\"", "\"") : "";
    }

    public static void main(String[] args) throws Exception {
        if (System.getenv("QUANTICDATA_API_KEY") == null) {
            System.err.println("set QUANTICDATA_API_KEY — https://app.quanticdata.io/register");
            System.exit(1);
        }

        String page = call("/scrape", "{\"url\":\"https://example.com\",\"format\":\"markdown\"}");
        String markdown = field(page, "markdown");
        System.out.println("--- markdown -------------------------------------------------");
        System.out.println(markdown.substring(0, Math.min(400, markdown.length())));
        System.out.println("title: " + field(page, "title"));

        String serp = call("/serp", "{\"query\":\"web scraping api\",\"country\":\"us\",\"num\":5}");
        System.out.println();
        System.out.println("--- organic (raw JSON, first 600 chars) ----------------------");
        System.out.println(serp.substring(0, Math.min(600, serp.length())));
    }
}
