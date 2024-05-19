package memorystorage

import (
	"context"
	"sync"

	"github.com/google/uuid"
	"github.com/trb1maker/otus_golang_home_work/hw12_13_14_15_calendar/internal/storage"
)

type Storage struct {
	mu    sync.RWMutex
	store map[string]storage.Event
}

func New() *Storage {
	return &Storage{}
}

func (s *Storage) InsertOne(_ context.Context, event storage.Event) (eventID string, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for {
		event.ID = uuid.New().String()
		if _, ok := s.store[event.ID]; !ok {
			break
		}
	}

	s.store[event.ID] = event
	return event.ID, nil
}

func (s *Storage) SelectOne(_ context.Context, eventID string) (storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	e, ok := s.store[eventID]
	if !ok {
		return storage.Event{}, storage.ErrNotFound
	}
	return e, nil
}

func (s *Storage) UpdateOne(_ context.Context, update storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.store[update.ID]; !ok {
		return storage.ErrNotFound
	}

	s.store[update.ID] = update
	return nil
}

func (s *Storage) DeleteOne(_ context.Context, eventID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.store[eventID]; !ok {
		return storage.ErrNotFound
	}

	delete(s.store, eventID)
	return nil
}
