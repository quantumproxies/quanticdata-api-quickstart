// QuanticData API quickstart — Go, standard library only.
//
//	QUANTICDATA_API_KEY=qd_live_... go run quickstart.go
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

const base = "https://api.quanticdata.io/v1"

type envelope struct {
	Type    string          `json:"type"`
	Message string          `json:"message"`
	Payload json.RawMessage `json:"payload"`
}

var client = &http.Client{Timeout: 120 * time.Second}

// call POSTs to the API and returns the payload, unwrapped from the envelope.
func call(path string, body any, out any) error {
	key := os.Getenv("QUANTICDATA_API_KEY")
	buf, _ := json.Marshal(body)
	req, err := http.NewRequest(http.MethodPost, base+path, bytes.NewReader(buf))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+key)
	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	raw, _ := io.ReadAll(res.Body)
	var env envelope
	if err := json.Unmarshal(raw, &env); err != nil {
		return fmt.Errorf("%s: bad JSON (HTTP %d)", path, res.StatusCode)
	}
	if env.Type == "error" || res.StatusCode >= 400 {
		return fmt.Errorf("%s failed (%d): %s", path, res.StatusCode, env.Message)
	}
	return json.Unmarshal(env.Payload, out)
}

type scrapePayload struct {
	Markdown string `json:"markdown"`
	Metadata struct {
		Title string `json:"title"`
	} `json:"metadata"`
	Usage struct {
		CostUSD float64 `json:"cost_usd"`
	} `json:"usage"`
}

type serpPayload struct {
	Organic []struct {
		Position int    `json:"position"`
		Title    string `json:"title"`
		Link     string `json:"link"`
	} `json:"organic"`
}

func main() {
	if os.Getenv("QUANTICDATA_API_KEY") == "" {
		fmt.Fprintln(os.Stderr, "set QUANTICDATA_API_KEY — https://app.quanticdata.io/register")
		os.Exit(1)
	}

	var page scrapePayload
	if err := call("/scrape", map[string]any{"url": "https://example.com", "format": "markdown"}, &page); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println("--- markdown -------------------------------------------------")
	if len(page.Markdown) > 400 {
		fmt.Println(page.Markdown[:400])
	} else {
		fmt.Println(page.Markdown)
	}
	fmt.Println("title:", page.Metadata.Title)
	fmt.Printf("cost : %v\n", page.Usage.CostUSD)

	var serp serpPayload
	if err := call("/serp", map[string]any{"query": "web scraping api", "country": "us", "num": 5}, &serp); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println("\n--- organic --------------------------------------------------")
	for _, r := range serp.Organic {
		fmt.Printf("%2d. %s\n    %s\n", r.Position, r.Title, r.Link)
	}
}
