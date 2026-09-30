// Example of calling the Validation API from another Go project using only
// net/http. Start the API first (go run main.go), then in another terminal:
//
//	go run ./examples/http_client
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

type validateResponse struct {
	Success bool   `json:"success"`
	Valid   bool   `json:"valid"`
	Type    string `json:"type"`
	Message string `json:"message"`
}

func main() {
	// 1. Create the JSON request
	requestBody, err := json.Marshal(map[string]string{
		"value": "ABCDE1234F",
	})
	if err != nil {
		fmt.Println("could not create request:", err)
		os.Exit(1)
	}

	req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/api/v1/validate/pan", bytes.NewReader(requestBody))
	if err != nil {
		fmt.Println("could not build request:", err)
		os.Exit(1)
	}

	// 3. Set Content-Type (and the API key, if the service requires one)
	req.Header.Set("Content-Type", "application/json")
	if key := os.Getenv("VALIDATION_API_KEY"); key != "" {
		req.Header.Set("X-API-Key", key)
	}

	// 2. Send the POST request (always set a timeout)
	httpClient := &http.Client{Timeout: 5 * time.Second}
	resp, err := httpClient.Do(req)
	if err != nil {
		// 7. Handle errors: service down, network problem, timeout...
		fmt.Println("request failed:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	// 4 + 5. Read the response and decode the JSON
	var result validateResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		fmt.Println("could not decode response:", err)
		os.Exit(1)
	}

	// Handle API-level errors (bad request, unauthorized, unsupported type...)
	if !result.Success {
		fmt.Printf("service returned an error (HTTP %d): %s\n", resp.StatusCode, result.Message)
		os.Exit(1)
	}

	// 6. Check the "valid" field
	if result.Valid {
		fmt.Println("PAN is valid:", result.Message)
	} else {
		fmt.Println("PAN is NOT valid:", result.Message)
	}
}
