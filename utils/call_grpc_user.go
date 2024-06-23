package utils

import (
	"context"
	"fmt"

	pos "github.com/Andrewalifb/alpha-pos-system-sales-service/api/proto"
	"google.golang.org/grpc"
)

// GetPosUser makes a gRPC request to the specified service and returns the response
func GetPosUserById(conn *grpc.ClientConn, id string, jwtPayload *pos.JWTPayload) (*pos.ReadPosUserResponse, error) {

	// Create a new PosUserService client
	client := pos.NewPosUserServiceClient(conn)

	// Prepare the request
	req := &pos.ReadPosUserRequest{
		UserId:     id,
		JwtPayload: jwtPayload,
	}

	// Call the gRPC method
	resp, err := client.ReadPosUser(context.Background(), req)
	if err != nil {
		return nil, fmt.Errorf("failed to call gRPC method: %w", err)
	}

	return resp, nil
}
