package dto

import pb "github.com/Andrewalifb/alpha-pos-system-sales-service/api/proto"

// // JSONTime is used to handle the protobuf timestamp fields in Go using time.Time
// type JSONTime struct {
// 	time.Time
// }

// // UnmarshalJSON converts JSON timestamp into the JSONTime format
// func (jt *JSONTime) UnmarshalJSON(b []byte) error {
// 	var timestamp struct {
// 		Seconds int64 `json:"seconds"`
// 		Nanos   int64 `json:"nanos"`
// 	}
// 	if err := json.Unmarshal(b, &timestamp); err != nil {
// 		return err
// 	}
// 	jt.Time = time.Unix(timestamp.Seconds, timestamp.Nanos)
// 	return nil
// }

// // MarshalJSON converts JSONTime back into JSON format
// func (jt JSONTime) MarshalJSON() ([]byte, error) {
// 	timestamp := struct {
// 		Seconds int64 `json:"seconds"`
// 		Nanos   int64 `json:"nanos"`
// 	}{
// 		Seconds: jt.Unix(),
// 		Nanos:   int64(jt.Nanosecond()),
// 	}
// 	return json.Marshal(timestamp)
// }

// UpdatePosProductRequest is the request structure for updating product information
type UpdateProductApiRequest struct {
	PosProduct *PosProduct    `json:"pos_product"`
	JwtPayload *pb.JWTPayload `json:"jwt_payload"`
	JwtToken   string         `json:"jwt_token"`
}

// SuccessResponse is used to parse the successful JSON response from the API
type UpdateProductApiSuccessResponse struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Data    *PosProduct `json:"data"`
}
