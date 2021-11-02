package event

import (
	"context"

	pb "github.com/nhatminhk63j/uetvoting/pb/event/v1"
	"github.com/nhatminhk63j/uetvoting/pkg/event"
)

type ServiceServer struct {
	*pb.UnimplementedEventServiceServer
	eventSvc event.Service
}

// NewServiceServer ...
func NewServiceServer(eventSvc event.Service) *ServiceServer {
	return &ServiceServer{
		eventSvc: eventSvc,
	}
}

// UpsertEvent ...
func (s *ServiceServer) UpsertEvent(ctx context.Context, in *pb.UpsertEventRequest) (*pb.UpsertEventResponse, error) {
	upsertInfo := event.ToEventUpsertInfo(in)
	eventID, err := s.eventSvc.UpsertEvent(ctx, upsertInfo)
	if err != nil {
		return nil, err
	}
	return &pb.UpsertEventResponse{
		Id: int32(eventID),
	}, nil
}

// GetEventByID ...
func (s *ServiceServer) GetEventByID(ctx context.Context, in *pb.GetEventByIDRequest) (*pb.GetEventByIDResponse, error) {
	eventDetail, err := s.eventSvc.GetEventDetail(ctx, int(in.Id))
	if err != nil {
		return nil, err
	}

	response, err := eventDetail.ToProtoStruct()
	if err != nil {
		return nil, err
	}
	return response, nil
}
