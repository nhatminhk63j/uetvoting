package authorization

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// NoPermissionError ...
type NoPermissionError struct {
}

// Error ...
func (e NoPermissionError) Error() string {
	return "no permission to access this RPC"
}

// GRPCStatus contains grpc status details.
func (e NoPermissionError) GRPCStatus() *status.Status {
	return status.New(codes.PermissionDenied, e.Error())
}
