package scheduler

import (
	"context"

	"github.com/mailru/easyjson"
	"github.com/trb1maker/otus_golang_home_work/hw12_13_14_15_calendar/internal/storage"
)

type Storage interface {
	SelectEventsToNotify(ctx context.Context) ([]storage.Notify, error)
	DeleteOldEvents(ctx context.Context) error
}

type Queue interface {
	Send(ctx context.Context, body []byte) error
}

type Scheduler struct {
	store Storage
	queue Queue
}

func New(store Storage, queue Queue) *Scheduler {
	return &Scheduler{
		store: store,
		queue: queue,
	}
}

func (s *Scheduler) Notify(ctx context.Context) error {
	events, err := s.store.SelectEventsToNotify(ctx)
	if err != nil {
		return err
	}
	for _, event := range events {
		body, err := easyjson.Marshal(event)
		if err != nil {
			return err
		}
		if err := s.queue.Send(ctx, body); err != nil {
			return err
		}
	}
	return nil
}

func (s *Scheduler) Clean(ctx context.Context) error {
	return s.store.DeleteOldEvents(ctx)
}
