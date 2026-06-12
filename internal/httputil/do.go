package httputil

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Same as http.Do but this builds the entire request - with body and headers
func Do(client *http.Client, url string, method string, headers map[string]string, body any) ([]byte, error) {
	jsonRequestData, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(
		method, url, bytes.NewBuffer(jsonRequestData),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error request to [%s]: %w", url, err)
	}

	// _ is actually an error but a close error is almost never actionable
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %s", resp.Status)
	}

	byteResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response bytes: %w", err)
	}

	return byteResp, nil
}
