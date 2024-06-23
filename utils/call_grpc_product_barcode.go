package utils

import (
	"context"
	"fmt"

	pos "github.com/Andrewalifb/alpha-pos-system-sales-service/api/proto"
	"google.golang.org/grpc"
)

// GetPosProduct makes a gRPC request to the specified service and returns the response
func GetPosProducByBarcode(conn *grpc.ClientConn, id string, jwtPayload *pos.JWTPayload, token string) (*pos.ReadPosProductByBarcodeResponse, error) {

	// Create a new PosProductService client
	client := pos.NewPosProductServiceClient(conn)

	// Prepare the request
	req := &pos.ReadPosProductByBarcodeRequest{
		ProductBarcodeId: id,
		JwtPayload:       jwtPayload,
		JwtToken:         token,
	}

	// Call the gRPC method
	resp, err := client.ReadPosProductByBarcode(context.Background(), req)
	if err != nil {
		return nil, fmt.Errorf("failed to call gRPC method: %w", err)
	}

	return resp, nil
}
