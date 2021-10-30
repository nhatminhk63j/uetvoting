package handler

import (
	"context"

	"github.com/nhatminhk63j/uetvoting/pb/health/v1"
)

type HealthServiceServer struct {
	*health.UnimplementedHealthCheckServiceServer
}

// Liveness ...
func (*HealthServiceServer) Liveness(context.Context, *health.LivenessRequest) (*health.LivenessResponse, error) {
	return &health.LivenessResponse{
		Content: "ok",
	}, nil
}

// Readiness ...
func (*HealthServiceServer) Readiness(context.Context, *health.ReadinessRequest) (*health.ReadinessResponse, error) {
	return &health.ReadinessResponse{
		Content: "ok",
	}, nil
}
