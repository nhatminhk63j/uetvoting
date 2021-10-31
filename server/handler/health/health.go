package health

import (
	"context"

	"github.com/nhatminhk63j/uetvoting/pb/health/v1"
)

type ServiceServer struct {
	*health.UnimplementedHealthCheckServiceServer
}

// Liveness ...
func (*ServiceServer) Liveness(context.Context, *health.LivenessRequest) (*health.LivenessResponse, error) {
	return &health.LivenessResponse{
		Content: "ok",
	}, nil
}

// Readiness ...
func (*ServiceServer) Readiness(context.Context, *health.ReadinessRequest) (*health.ReadinessResponse, error) {
	return &health.ReadinessResponse{
		Content: "ok",
	}, nil
}
