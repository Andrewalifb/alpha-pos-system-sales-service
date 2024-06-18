package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	pb "github.com/Andrewalifb/alpha-pos-system-sales-service/api/proto"
	"github.com/Andrewalifb/alpha-pos-system-sales-service/dto"
	"io"
	"net/http"
	"time"
)

// UpdatePosProduct sends a PUT request to update a POS product by its ID
func UpdatePosProduct(product *dto.PosProduct, jwtPayload *pb.JWTPayload, token string, productId string) (*dto.PosProduct, error) {
	// Construct the URL for the request
	url := fmt.Sprintf("http://localhost:8082/api/v1/products/pos_product/%s", productId)

	// Prepare the request body from the product data
	requestBody := dto.UpdateProductApiRequest{
		PosProduct: product,
		JwtPayload: jwtPayload,
		JwtToken:   token,
	}

	// Marshal the request body into JSON
	bodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	// Create the HTTP PUT request
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(bodyBytes))
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
	var successResp dto.UpdateProductApiSuccessResponse
	if err := json.Unmarshal(respBody, &successResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal success response: %w", err)
	}

	// Debugging line - Normally this should be handled by a logger in production code
	fmt.Printf("UPDATE POS PRODUCT RESPONSE: %+v\n", successResp)

	// Check if the Data part of the response is nil
	if successResp.Data == nil {
		return nil, fmt.Errorf("pos product data is missing in the response")
	}

	// Return the updated PosProduct data from the response
	return successResp.Data, nil
}
