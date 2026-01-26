package client

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

func PostModelApi(url string, reqJson string, apiKey string) (string, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBufferString(reqJson))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("remote returned %d: %s", resp.StatusCode, string(body))
	}

	return string(body), nil
}
