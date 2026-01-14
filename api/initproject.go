package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type InitProjectRequest struct {
	Name string `json:"name"`
}

type InitProjectResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func CreateProjectOnServer(projectName string) (string, error) {
	apiUrl := "http://localhost:8000/cli/init"

	reqBody := InitProjectRequest{
		Name: projectName,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := http.Post(
		apiUrl,
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return "", fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("server returned error: %s (status: %d)", string(body), resp.StatusCode)
	}

	var response InitProjectResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	return response.ID, nil
}
