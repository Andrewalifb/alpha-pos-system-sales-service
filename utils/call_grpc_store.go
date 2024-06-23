package utils

import (
	"context"
	"fmt"

	pos "github.com/Andrewalifb/alpha-pos-system-sales-service/api/proto"
	"google.golang.org/grpc"
)

// GetPosStore makes a gRPC request to the specified service and returns the response
func GetPosStoreById(conn *grpc.ClientConn, id string, jwtPayload *pos.JWTPayload) (*pos.ReadPosStoreResponse, error) {

	// Create a new PosStoreService client
	client := pos.NewPosStoreServiceClient(conn)

	// Prepare the request
	req := &pos.ReadPosStoreRequest{
		StoreId:    id,
		JwtPayload: jwtPayload,
	}

	// Call the gRPC method
	resp, err := client.ReadPosStore(context.Background(), req)
	if err != nil {
		return nil, fmt.Errorf("failed to call gRPC method: %w", err)
	}

	return resp, nil
}
