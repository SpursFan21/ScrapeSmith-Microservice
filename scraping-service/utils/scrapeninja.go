package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func ScrapeWithScrapeNinja(url string) ([]byte, error) {
	apiURL := fmt.Sprintf("https://%s/scrape-js", os.Getenv("SCRAPENINJA_API_HOST"))

	requestBody, _ := json.Marshal(map[string]interface{}{
		"url":      url,
		"geo":      "us",
		"retryNum": 1,
	})

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-RapidAPI-Key", os.Getenv("SCRAPENINJA_API_KEY"))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return responseData, nil
}
