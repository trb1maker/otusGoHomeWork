package internalgrpc

import (
	"context"
	"net"
	"strconv"
	"time"

	"github.com/trb1maker/otus_golang_home_work/hw12_13_14_15_calendar/internal/server/grpc/api"
	"github.com/trb1maker/otus_golang_home_work/hw12_13_14_15_calendar/internal/storage"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Application interface {
	CreateEvent(ctx context.Context, userID string, title string, startTime time.Time,
		endTime time.Time, description string, notify time.Duration) (string, error)
	GetEvent(ctx context.Context, eventID string) (storage.Event, error)
	UpdateEvent(ctx context.Context, event storage.Event) error
	DeleteEvent(ctx context.Context, eventID string) error
	GetNextEvent(ctx context.Context, ownerID string) (storage.Event, error)
	GetAllEvents(ctx context.Context, ownerID string) ([]storage.Event, error)
	GetEventsFromRange(ctx context.Context, ownerID string, from time.Time, to time.Time) ([]storage.Event, error)
	GetEventsCurrentDay(ctx context.Context, ownerID string) ([]storage.Event, error)
	GetEventsCurrentWeek(ctx context.Context, ownerID string) ([]storage.Event, error)
	GetEventsCurrentMonth(ctx context.Context, ownerID string) ([]storage.Event, error)
}

type Server struct {
	api.UnimplementedEventServiceServer
	app     Application
	srv     *grpc.Server
	address string
}

func NewServer(app Application, host string, port int) *Server {
	return &Server{
		app:     app,
		address: net.JoinHostPort(host, strconv.Itoa(port)),
	}
}

func (s *Server) Start() error {
	lis, err := net.Listen("tcp", s.address)
	if err != nil {
		return err
	}
	defer lis.Close()

	s.srv = grpc.NewServer(grpc.UnaryInterceptor(logginInterceptor))

	api.RegisterEventServiceServer(s.srv, s)

	return s.srv.Serve(lis)
}

func (s *Server) Stop() {
	s.srv.Stop()
}

func (s *Server) NewEvent(ctx context.Context, e *api.Event) (*api.EventIdResponse, error) {
	eventID, err := s.app.CreateEvent(
		ctx,
		e.GetOwner(),
		e.GetTitle(),
		e.StartTime.AsTime(),
		e.EndTime.AsTime(),
		e.GetDescription(),
		e.Notify.AsDuration(),
	)
	if err != nil {
		return nil, err
	}
	return &api.EventIdResponse{EventId: eventID}, nil
}

func (s *Server) GetEvent(ctx context.Context, r *api.EventRequest) (*api.EventResponse, error) {
	event, err := s.app.GetEvent(ctx, r.GetEventId())
	if err != nil {
		return nil, err
	}
	return &api.EventResponse{Events: []*api.Event{{
		Id:          event.ID,
		Title:       event.Title,
		StartTime:   timestamppb.New(event.StartTime),
		EndTime:     timestamppb.New(event.EndTime),
		Owner:       event.OwnerID,
		Description: event.Description,
		Notify:      durationpb.New(event.Notify),
	}}}, nil
}

func (s *Server) UpdateEvent(ctx context.Context, e *api.Event) (*emptypb.Empty, error) {
	if err := s.app.UpdateEvent(ctx, storage.Event{
		ID:          e.GetId(),
		Title:       e.GetTitle(),
		StartTime:   e.StartTime.AsTime(),
		EndTime:     e.EndTime.AsTime(),
		OwnerID:     e.GetOwner(),
		Description: e.GetDescription(),
		Notify:      e.Notify.AsDuration(),
	}); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) DeleteEvent(ctx context.Context, r *api.EventRequest) (*emptypb.Empty, error) {
	if err := s.app.DeleteEvent(ctx, r.GetEventId()); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func transform(ee ...storage.Event) *api.EventResponse {
	dto := &api.EventResponse{Events: make([]*api.Event, 0, len(ee))}
	for i := 0; i < len(ee); i++ {
		dto.Events = append(dto.Events, &api.Event{
			Id:          ee[i].ID,
			Title:       ee[i].Title,
			StartTime:   timestamppb.New(ee[i].StartTime),
			EndTime:     timestamppb.New(ee[i].EndTime),
			Owner:       ee[i].OwnerID,
			Description: ee[i].Description,
			Notify:      durationpb.New(ee[i].Notify),
		})
	}
	return dto
}

func (s *Server) All(ctx context.Context, r *api.UserRequest) (*api.EventResponse, error) {
	events, err := s.app.GetAllEvents(ctx, r.GetUserId())
	if err != nil {
		return nil, err
	}
	return transform(events...), nil
}

func (s *Server) Next(ctx context.Context, r *api.UserRequest) (*api.EventResponse, error) {
	event, err := s.app.GetNextEvent(ctx, r.GetUserId())
	if err != nil {
		return nil, err
	}
	return transform(event), nil
}

func (s *Server) Day(ctx context.Context, r *api.UserRequest) (*api.EventResponse, error) {
	events, err := s.app.GetEventsCurrentDay(ctx, r.GetUserId())
	if err != nil {
		return nil, err
	}
	return transform(events...), nil
}

func (s *Server) Week(ctx context.Context, r *api.UserRequest) (*api.EventResponse, error) {
	events, err := s.app.GetEventsCurrentWeek(ctx, r.GetUserId())
	if err != nil {
		return nil, err
	}
	return transform(events...), nil
}

func (s *Server) Month(ctx context.Context, r *api.UserRequest) (*api.EventResponse, error) {
	events, err := s.app.GetEventsCurrentMonth(ctx, r.GetUserId())
	if err != nil {
		return nil, err
	}
	return transform(events...), nil
}
