package utils

import (
	"context"
	"fmt"

	pos "github.com/Andrewalifb/alpha-pos-system-sales-service/api/proto"
	"google.golang.org/grpc"
)

// CreatePosInventoryHistory makes a gRPC request to the specified service and returns the response
func CreateNewPosInventoryHistory(conn *grpc.ClientConn, inventory *pos.PosInventoryHistory, jwtPayload *pos.JWTPayload, token string) (*pos.CreatePosInventoryHistoryResponse, error) {

	// Create a new PosInventoryHistoryService client
	client := pos.NewPosInventoryHistoryServiceClient(conn)

	// Prepare the request
	req := &pos.CreatePosInventoryHistoryRequest{
		PosInventoryHistory: inventory,
		JwtPayload:          jwtPayload,
		JwtToken:            token,
	}

	// Call the gRPC method
	resp, err := client.CreatePosInventoryHistory(context.Background(), req)
	if err != nil {
		return nil, fmt.Errorf("failed to call gRPC method: %w", err)
	}

	return resp, nil
}
