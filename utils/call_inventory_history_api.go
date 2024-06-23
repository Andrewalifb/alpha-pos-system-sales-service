package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	pb "github.com/Andrewalifb/alpha-pos-system-sales-service/api/proto"
	"github.com/Andrewalifb/alpha-pos-system-sales-service/dto"
)

// CreatePosInventoryHistory sends a POST request to create a new POS inventory history record
func CreatePosInventoryHistory(history *dto.PosInventoryHistory, jwtPayload *pb.JWTPayload, token string) (*dto.PosInventoryHistory, error) {
	// URL of the API endpoint
	url := "http://localhost:8082/api/v1/inventory-histories/pos_inventory_history"

	// Prepare the request body from the history data
	requestBody := dto.CreatePosInventoryHistoryRequest{
		PosInventoryHistory: history,
		JwtPayload:          jwtPayload,
		JwtToken:            token,
	}

	// Marshal the request body into JSON
	bodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	// Create the HTTP POST request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Set the necessary headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	// Initialize the HTTP client with a timeout
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer resp.Body.Close()

	// Check if the HTTP response status is not 200 OK
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 response status: %d %s", resp.StatusCode, http.StatusText(resp.StatusCode))
	}

	// Read the response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Unmarshal the JSON response into the SuccessResponse DTO
	var successResp dto.CreateInventorySuccessResponse
	if err := json.Unmarshal(respBody, &successResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal success response: %w", err)
	}

	// Debugging line - Normally this should be handled by a logger in production code
	fmt.Printf("CREATE INVENTORY HISTORY RESPONSE: %+v\n", successResp)

	// Check if the Data part of the response is nil or if PosInventoryHistory is nil
	if successResp.Data == nil || successResp.Data.PosInventoryHistory == nil {
		return nil, fmt.Errorf("pos inventory history data is missing in the response")
	}

	// Return the created PosInventoryHistory data from the response
	return successResp.Data.PosInventoryHistory, nil
}
