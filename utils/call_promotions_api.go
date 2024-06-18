package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Andrewalifb/alpha-pos-system-sales-service/dto"
)

// GetPosPromotion fetches promotion details by ID and JWT token
func GetPosPromotion(id string, token string) (*dto.PosPromotion, error) {
	// Construct the URL for the request
	url := fmt.Sprintf("http://localhost:8082/api/v1/promotions/pos_promotion/%s", id)

	// Create a new request using http
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to make HTTP request: %w", err)
	}

	// Add headers
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+token)

	// Send the request using a new HTTP client with a timeout
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Check if the response status code is not 200
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 response status: %d %s", resp.StatusCode, http.StatusText(resp.StatusCode))
	}

	// Read the response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Unmarshal the response body into the SuccessResponse DTO
	var successResp dto.ReadPromotionApiSuccessResponse
	if err := json.Unmarshal(respBody, &successResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal success response: %w", err)
	}

	fmt.Printf("PROMOTION DATA: %+v\n", successResp.Data.PosPromotion)

	// Check if the PosPromotion is nil to avoid dereferencing a nil pointer
	if successResp.Data == nil || successResp.Data.PosPromotion == nil {
		return nil, fmt.Errorf("promotion data is missing in the response")
	}

	return successResp.Data.PosPromotion, nil
}

// GetPosPromotionByProductID fetches promotion details by Product ID and JWT token
func GetPosPromotionByProductID(productID string, token string) (*dto.PosPromotion, error) {
	// Construct the URL for the request
	url := fmt.Sprintf("http://localhost:8082/api/v1/promotions/pos_promotion/by_product/%s", productID)

	// Create a new request using http
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to make HTTP request: %w", err)
	}

	// Add headers
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+token)

	// Send the request using a new HTTP client with a timeout
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Check if the response status code is not 200
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 response status: %d %s", resp.StatusCode, http.StatusText(resp.StatusCode))
	}

	// Read the response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Unmarshal the response body into the SuccessResponse DTO
	var successResp dto.ReadPromotionApiSuccessResponse
	if err := json.Unmarshal(respBody, &successResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal success response: %w", err)
	}

	fmt.Printf("PROMOTION DATA: %+v\n", successResp.Data.PosPromotion)

	// Check if the PosPromotion is nil to avoid dereferencing a nil pointer
	if successResp.Data == nil || successResp.Data.PosPromotion == nil {
		return nil, fmt.Errorf("promotion data is missing in the response")
	}

	return successResp.Data.PosPromotion, nil
}
