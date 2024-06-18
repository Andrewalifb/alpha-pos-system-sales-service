package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Andrewalifb/alpha-pos-system-sales-service/dto"
)

// GetPosProduct makes a GET request to the specified API endpoint and returns the response
func GetPosProduct(id string, token string) (*dto.ReadProductApiResponse, error) {
	// Create a new request using http
	req, err := http.NewRequest("GET", "http://localhost:8082/api/v1/products/pos_product_barcode/"+id, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to make HTTP request: %w", err)
	}

	// Add headers
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+token)

	// Send the request using a new http Client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}
	defer resp.Body.Close()

	// Read the response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Unmarshal the response body to the response DTO
	var respDto dto.ReadProductApiResponse
	err = json.Unmarshal(respBody, &respDto)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response DTO: %w", err)
	}
	fmt.Println("PRODUCT DATA:", respDto.Data.PosProduct)
	return &respDto, nil
}
