package authentication

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// MissingMetadataError ...
type MissingMetadataError struct {
}

// Error ...
func (e MissingMetadataError) Error() string {
	return "metadata is not provided"
}

// GRPCStatus contains grpc status details.
func (e MissingMetadataError) GRPCStatus() *status.Status {
	return status.New(codes.Unauthenticated, e.Error())
}

// MissingAccessTokenError ...
type MissingAccessTokenError struct {
}

// Error ...
func (e MissingAccessTokenError) Error() string {
	return "authorization token is not provided"
}

// GRPCStatus contains grpc status details.
func (e MissingAccessTokenError) GRPCStatus() *status.Status {
	return status.New(codes.Unauthenticated, e.Error())
}

// InvalidAccessTokenError ...
type InvalidAccessTokenError struct {
}

// Error ...
func (e InvalidAccessTokenError) Error() string {
	return "access token is invalid"
}

// GRPCStatus contains grpc status details.
func (e InvalidAccessTokenError) GRPCStatus() *status.Status {
	return status.New(codes.Unauthenticated, e.Error())
}
