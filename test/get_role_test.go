package utils

// import (
// 	"net/http"
// 	"net/http/httptest"
// 	"os"
// 	"testing"

// 	pb "github.com/Andrewalifb/alpha-pos-system-sales-service/api/proto"
// 	"github.com/Andrewalifb/alpha-pos-system-sales-service/utils"
// 	"github.com/stretchr/testify/assert"
// )

// func TestGetPosRole(t *testing.T) {
// 	// Create a mock server
// 	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

// 		assert.Equal(t, req.URL.String(), "/api/v1/roles/pos_role/test_role_id")

// 		rw.Write([]byte(`{"role_id": "test_role_id", "role_name": "test_role_name"}`))
// 	}))
// 	// Close the server when test finishes
// 	defer server.Close()

// 	os.Setenv("SERVER_URL", server.URL)

// 	protoPayload := &pb.JWTPayload{
// 		Name: "test_name",
// 		Role: "test_role",
// 	}
// 	token := "test_token"

// 	resp, err := utils.GetPosRole("test_role_id", protoPayload, token)

// 	assert.Nil(t, err)

// 	assert.Equal(t, "test_role_id", resp.RoleID)
// 	assert.Equal(t, "test_role_name", resp.RoleName)
// }
