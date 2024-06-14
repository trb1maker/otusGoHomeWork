package sender

import (
	"context"
	"log/slog"
	"time"

	"github.com/mailru/easyjson"
	"github.com/trb1maker/otus_golang_home_work/hw12_13_14_15_calendar/internal/storage"
)

type Queue interface {
	Consume(ctx context.Context) (<-chan []byte, error)
}

func New(queue Queue) *Sender {
	return &Sender{queue: queue}
}

type Sender struct {
	queue Queue
}

func (s *Sender) Schedule(ctx context.Context) error {
	events, err := s.queue.Consume(ctx)
	if err != nil {
		return err
	}
	var event storage.Notify
	for data := range events {
		if err := easyjson.Unmarshal(data, &event); err != nil {
			return err
		}
		go s.Notify(ctx, event)
	}
	return nil
}

func (s *Sender) Notify(ctx context.Context, event storage.Notify) {
	time.Sleep(time.Until(event.NotifyTime))
	slog.InfoContext(
		ctx,
		"notification",
		slog.Group(
			"event",
			slog.String("userId", event.OwnerID),
			slog.String("eventId", event.ID),
			slog.String("title", event.Title),
			slog.String("time", event.StartTime.Format(time.TimeOnly)),
		),
	)
}
