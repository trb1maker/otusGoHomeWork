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
	SelectAllEvents(context.Context, string) ([]storage.Event, error)
	SelectNextEvent(context.Context, string) (storage.Event, error)
	SelectEventsBetweenDates(context.Context, string, time.Time, time.Time) ([]storage.Event, error)
	UpdateOne(context.Context, storage.Event) error
	DeleteOne(context.Context, string) error
}

func New(storage Storage) *App {
	return &App{
		store: storage,
	}
}

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

func (a *App) GetAllEvents(ctx context.Context, ownerID string) ([]storage.Event, error) {
	return a.store.SelectAllEvents(ctx, ownerID)
}

func (a *App) GetNextEvent(ctx context.Context, ownerID string) (storage.Event, error) {
	return a.store.SelectNextEvent(ctx, ownerID)
}

func (a *App) GetEventsFromRange(
	ctx context.Context,
	ownerID string,
	from time.Time,
	to time.Time,
) ([]storage.Event, error) {
	return a.store.SelectEventsBetweenDates(ctx, ownerID, from, to)
}

func (a *App) GetEventsCurrentDay(ctx context.Context, ownerID string) ([]storage.Event, error) {
	now := time.Now()

	return a.GetEventsFromRange(ctx, ownerID, startOfDay(now), endOfDay(now))
}

func (a *App) GetEventsCurrentWeek(ctx context.Context, ownerID string) ([]storage.Event, error) {
	now := time.Now()

	return a.GetEventsFromRange(ctx, ownerID, startOfWeek(now), endOfWeek(now))
}

func (a *App) GetEventsCurrentMonth(ctx context.Context, ownerID string) ([]storage.Event, error) {
	now := time.Now()

	return a.GetEventsFromRange(ctx, ownerID, startOfMonth(now), endOfMonth(now))
}
