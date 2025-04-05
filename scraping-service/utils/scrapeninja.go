//scraping-service\utils\scrapeninja.go
package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type ScrapeNinjaResponse struct {
	Info map[string]interface{} `json:"info"`
	Body string                 `json:"body"`
}

func ScrapeWithScrapeNinja(url string) (string, error) {
	apiURL := fmt.Sprintf("https://%s/scrape-js", os.Getenv("SCRAPENINJA_API_HOST"))

	requestBody, _ := json.Marshal(map[string]interface{}{
		"url":      url,
		"geo":      "us",
		"retryNum": 1,
	})

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-RapidAPI-Key", os.Getenv("SCRAPENINJA_API_KEY"))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var parsed ScrapeNinjaResponse
	if err := json.Unmarshal(bodyBytes, &parsed); err != nil {
		return "", err
	}

	return parsed.Body, nil
}
