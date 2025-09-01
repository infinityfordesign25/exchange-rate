package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

const api = "https://open.er-api.com/v6/latest/"

var iso4217Re = regexp.MustCompile(`^[A-Z]{3}$`)

type ratesResponse struct {
	Result             string             `json:"result"`
	TimeLastUpdateUnix int64              `json:"time_last_update_unix"`
	Rates              map[string]float64 `json:"rates"`
}

func fetchRates(ctx context.Context, client *http.Client, base string) (*ratesResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, api+base, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch URL: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed with status code: %d %s", resp.StatusCode, http.StatusText(resp.StatusCode))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var data ratesResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if data.Result != "success" {
		return nil, fmt.Errorf("failed to get \"success\": %s", data.Result)
	}

	return &data, nil
}

func main() {
	source := flag.String("source", "", "Source currency code, e.g. USD")
	target := flag.String("target", "", "Target currency code, e.g. EUR")
	amount := flag.Float64("amount", 1, "Amount in source currency")
	flag.Parse()

	*source = strings.ToUpper(*source)
	*target = strings.ToUpper(*target)

	if !iso4217Re.MatchString(*source) || !iso4217Re.MatchString(*target) {
		flag.Usage()
		os.Exit(2)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	data, err := fetchRates(ctx, client, *source)
	if err != nil {
		log.Fatalf("fetchRates: %v", err)
	}

	rate, ok := data.Rates[*target]
	if !ok {
		log.Fatalf("no rate for %s", *target)
	}

	lastUpdated := time.Unix(data.TimeLastUpdateUnix, 0)

	fmt.Printf("%.2f %s = %.2f %s\n", *amount, *source, (*amount)*rate, *target)
	fmt.Println("Last updated on:")
	fmt.Println("Local:", lastUpdated.Format(time.RFC3339))
	fmt.Println("UTC:  ", lastUpdated.UTC().Format(time.RFC3339))
}
