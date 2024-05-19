package app

import (
	"context"
	"time"

	"github.com/trb1maker/otus_golang_home_work/hw12_13_14_15_calendar/internal/storage"
)

type App struct {
	store Storage
}

type Storage interface {
	InsertOne(context.Context, storage.Event) (string, error)
	SelectOne(context.Context, string) (storage.Event, error)
	UpdateOne(context.Context, storage.Event) error
	DeleteOne(context.Context, string) error
}

func New(storage Storage) *App {
	return &App{
		store: storage,
	}
}

// TODO: типы аргументов должны соответствовать возможностям HTTP и GRPC,
// на текущий момент я реализую их в соответствии с типами Event,
// позже, возможно поменяю.

func (a *App) CreateEvent(
	ctx context.Context,
	userID string,
	title string,
	startTime time.Time,
	endTime time.Time,
	description string,
	notify time.Duration,
) (string, error) {
	e := storage.Event{
		Title:       title,
		StartTime:   startTime,
		EndTime:     endTime,
		Description: description,
		OwnerID:     userID,
		Notify:      notify,
	}
	return a.store.InsertOne(ctx, e)
}

func (a *App) GetEvent(ctx context.Context, eventID string) (storage.Event, error) {
	return a.store.SelectOne(ctx, eventID)
}

func (a *App) UpdateEvent(ctx context.Context, event storage.Event) error {
	return a.store.UpdateOne(ctx, event)
}

func (a *App) DeleteEvent(ctx context.Context, eventID string) error {
	return a.store.DeleteOne(ctx, eventID)
}

// TODO: Реализовать другие методы приложения после GRPC
