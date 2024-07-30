package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

const (
	listenAddr    = "localhost:8080"
	targetBaseURL = "https://api.example.com" // Replace with your target URL
)

type OriginalResponse struct {
	// Define the structure of the original response
	Data string `json:"data"`
}

type TransformedResponse struct {
	// Define the structure of the transformed response
	TransformedData string `json:"transformed_data"`
}

func main() {
	http.HandleFunc("/", proxyHandler)
	log.Printf("Starting proxy server on %s", listenAddr)
	log.Fatal(http.ListenAndServe(listenAddr, nil))
}

func proxyHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the original URL
	query := r.URL.Query().Get("query")
	if query == "" {
		http.Error(w, "Missing 'query' parameter", http.StatusBadRequest)
		return
	}

	// Construct the new URL with the query parameter
	targetURL, err := url.Parse(targetBaseURL)
	if err != nil {
		http.Error(w, "Error parsing target URL", http.StatusInternalServerError)
		return
	}
	targetURL.RawQuery = url.Values{"query": {query}}.Encode()

	// Send the request to the new URL
	resp, err := http.Get(targetURL.String())
	if err != nil {
		http.Error(w, "Error forwarding request", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Error reading response body", http.StatusInternalServerError)
		return
	}

	// Parse the original response
	var originalResp OriginalResponse
	if err := json.Unmarshal(body, &originalResp); err != nil {
		http.Error(w, "Error parsing original response", http.StatusInternalServerError)
		return
	}

	// Transform the response
	transformedResp := TransformedResponse{
		TransformedData: fmt.Sprintf("Transformed: %s", originalResp.Data),
	}

	// Send the transformed response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transformedResp)
}
