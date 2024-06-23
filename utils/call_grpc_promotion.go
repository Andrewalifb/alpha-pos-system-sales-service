package utils

import (
	"context"
	"fmt"

	pos "github.com/Andrewalifb/alpha-pos-system-sales-service/api/proto"
	"google.golang.org/grpc"
)

// GetPosPromotionByProductId makes a gRPC request to the specified service and returns the response
func GetPosPromotionByProductId(conn *grpc.ClientConn, id string, jwtPayload *pos.JWTPayload, token string) (*pos.ReadPosPromotionByProductIdResponse, error) {

	// Create a new PosPromotionService client
	client := pos.NewPosPromotionServiceClient(conn)

	// Prepare the request
	req := &pos.ReadPosPromotionByProductIdRequest{
		ProductId:  id,
		JwtPayload: jwtPayload,
		JwtToken:   token,
	}

	// Call the gRPC method
	resp, err := client.ReadPosPromotionByProductId(context.Background(), req)
	if err != nil {
		return nil, fmt.Errorf("failed to call gRPC method: %w", err)
	}

	return resp, nil
}
