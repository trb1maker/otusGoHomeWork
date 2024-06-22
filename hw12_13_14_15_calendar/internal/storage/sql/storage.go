package sqlstorage

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx"
	"github.com/trb1maker/otus_golang_home_work/hw12_13_14_15_calendar/internal/storage"
)

var ErrNotValidInterval = errors.New("interval equal zero")

type Storage struct {
	db             *pgx.Conn
	notifyInterval time.Duration
}

func New(host string, port uint16, dbName string, userName string, userPassword string) (*Storage, error) {
	config := &pgx.ConnConfig{
		Host:     host,
		Port:     port,
		Database: dbName,
		User:     userName,
		Password: userPassword,
	}
	db, err := pgx.Connect(*config)
	if err != nil {
		return nil, err
	}

	return &Storage{db: db}, nil
}

func (s *Storage) Connect(ctx context.Context) error {
	return s.db.Ping(ctx)
}

func (s *Storage) Close(_ context.Context) error {
	return s.db.Close()
}

func (s *Storage) InsertOne(ctx context.Context, event storage.Event) (eventID string, err error) {
	tx, err := s.db.BeginEx(ctx, nil)
	if err != nil {
		return "", fmt.Errorf("can't begin transaction %w", err)
	}
	defer tx.RollbackEx(ctx)

	result := tx.QueryRowEx(
		ctx,
		`
			insert into otus.events (title, start_time, end_time, description, user_id, notify)
			values ($1, $2, $3, $4, $5, $6)
			returning event_id
		`,
		nil,
		event.Title,
		event.StartTime.UTC(),
		event.EndTime.UTC(),
		event.Description,
		event.OwnerID,
		event.Notify,
	)
	if err := result.Scan(&eventID); err != nil {
		return "", fmt.Errorf("can't insert value %w", err)
	}
	return eventID, tx.CommitEx(ctx)
}

func (s *Storage) SelectOne(ctx context.Context, eventID string) (storage.Event, error) {
	tx, err := s.db.BeginEx(ctx, nil)
	if err != nil {
		return storage.Event{}, fmt.Errorf("can't begin transaction %w", err)
	}
	defer tx.RollbackEx(ctx)

	result := tx.QueryRowEx(
		ctx,
		`
			select title, start_time, end_time, description, user_id, notify
			from otus.events
			where event_id = $1
			  and is_deleted = false
		`,
		nil,
		eventID,
	)
	var e storage.Event
	if err := result.Scan(&e.Title, &e.StartTime, &e.EndTime, &e.Description, &e.OwnerID, &e.Notify); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return storage.Event{}, storage.ErrNotFound
		}
		return storage.Event{}, fmt.Errorf("can't scan result %w", err)
	}
	e.ID = eventID
	return e, tx.CommitEx(ctx)
}

func (s *Storage) UpdateOne(ctx context.Context, event storage.Event) error {
	tx, err := s.db.BeginEx(ctx, nil)
	if err != nil {
		return fmt.Errorf("can't begin transaction %w", err)
	}
	defer tx.RollbackEx(ctx)

	result, err := s.db.ExecEx(
		ctx,
		`
			update otus.events set
				title = $1,
				start_time = $2,
				end_time = $3,
				description = $4,
				notify = $5
			where event_id = $6
		`,
		nil,
		event.Title,
		event.StartTime.UTC(),
		event.EndTime.UTC(),
		event.Description,
		event.Notify,
		event.ID,
	)
	if err != nil {
		return fmt.Errorf("can't update event %w", err)
	}

	if result.RowsAffected() == 0 {
		return storage.ErrNotFound
	}

	return tx.CommitEx(ctx)
}

func (s *Storage) DeleteOne(ctx context.Context, eventID string) error {
	tx, err := s.db.BeginEx(ctx, nil)
	if err != nil {
		return fmt.Errorf("can't begin transaction %w", err)
	}
	defer tx.RollbackEx(ctx)

	result, err := tx.ExecEx(
		ctx,
		`
			update otus.events
			set is_deleted = true
			where event_id = $1
		`,
		nil,
		eventID,
	)
	if err != nil {
		return fmt.Errorf("can't delete event %w", err)
	}

	if result.RowsAffected() == 0 {
		return storage.ErrNotFound
	}

	return tx.CommitEx(ctx)
}

func (s *Storage) SelectAllEvents(ctx context.Context, userID string) ([]storage.Event, error) {
	rows, err := s.db.QueryEx(
		ctx,
		`
			select event_id, title, start_time, end_time, description, notify
			from otus.events
			where is_deleted = false 
			  and user_id = $1
		`,
		nil,
		userID,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, storage.ErrNoData
		}
		return nil, err
	}

	var (
		eventID, title, description string
		startTime, endTime          time.Time
		notify                      time.Duration

		events []storage.Event
	)

	for rows.Next() {
		if err := rows.Scan(&eventID, &title, &startTime, &endTime, &description, &notify); err != nil {
			return nil, err
		}
		events = append(events, storage.Event{
			ID:          eventID,
			Title:       title,
			StartTime:   startTime,
			EndTime:     endTime,
			Description: description,
			Notify:      notify,
			OwnerID:     userID,
		})
	}

	return events, nil
}

func (s *Storage) SelectNextEvent(ctx context.Context, userID string) (storage.Event, error) {
	result := s.db.QueryRowEx(
		ctx,
		`
			select event_id, title, start_time, end_time, description, notify
			from otus.events
			where is_deleted = false
			  and user_id = $1
			  and start_time >= current_timestamp
			order by start_time
			limit 1
		`,
		nil,
		userID,
	)

	var (
		eventID, title, description string
		startTime, endTime          time.Time
		notify                      time.Duration
	)
	if err := result.Scan(&eventID, &title, &startTime, &endTime, &description, &notify); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return storage.Event{}, storage.ErrNoData
		}
		return storage.Event{}, err
	}

	return storage.Event{
		ID:          eventID,
		Title:       title,
		StartTime:   startTime,
		EndTime:     endTime,
		Description: description,
		Notify:      notify,
		OwnerID:     userID,
	}, nil
}

func (s *Storage) SelectEventsBetweenDates(
	ctx context.Context,
	userID string,
	from time.Time,
	to time.Time,
) ([]storage.Event, error) {
	rows, err := s.db.QueryEx(
		ctx,
		`
			select event_id, title, start_time, end_time, description, notify
			from otus.events
			where is_deleted = false
			  and user_id = $1
			  and start_time >= $2
			  and start_time <= $3
		`,
		nil,
		userID,
		from,
		to,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, storage.ErrNoData
		}
		return nil, err
	}

	var (
		eventID, title, description string
		startTime, endTime          time.Time
		notify                      time.Duration

		events []storage.Event
	)

	for rows.Next() {
		if err := rows.Scan(&eventID, &title, &startTime, &endTime, &description, &notify); err != nil {
			return nil, err
		}
		events = append(events, storage.Event{
			ID:          eventID,
			Title:       title,
			StartTime:   startTime,
			EndTime:     endTime,
			Description: description,
			Notify:      notify,
			OwnerID:     userID,
		})
	}

	return events, nil
}

func (s *Storage) RegisterUser(ctx context.Context) (userID string, err error) {
	userID = uuid.New().String()
	if _, err := s.db.ExecEx(ctx, "insert into otus.users (user_id) values ($1)", nil, userID); err != nil {
		return "", err
	}
	return userID, err
}

func (s *Storage) SetInterval(interval time.Duration) error {
	if interval == 0 {
		return ErrNotValidInterval
	}
	s.notifyInterval = interval
	return nil
}

func (s *Storage) SelectEventsToNotify(ctx context.Context) ([]storage.Notify, error) {
	rows, err := s.db.QueryEx(
		ctx,
		`
			select
				user_id,
				event_id,
				title,
				start_time,
				start_time - "notify" notify_time
			from otus.events 
			where "notify" is not null
			and not is_deleted
			and start_time - "notify" between current_timestamp and current_timestamp + $1
		`,
		nil,
		s.notifyInterval,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, storage.ErrNoData
		}
		return nil, err
	}

	var (
		userID, eventID, title string
		startTime, notifyTime  time.Time

		notify []storage.Notify
	)
	for rows.Next() {
		if err := rows.Scan(&userID, &eventID, &title, &startTime, &notifyTime); err != nil {
			return nil, err
		}
		notify = append(notify, storage.Notify{
			OwnerID:    userID,
			ID:         eventID,
			Title:      title,
			StartTime:  startTime,
			NotifyTime: notifyTime,
		})
	}

	return notify, nil
}

func (s *Storage) DeleteOldEvents(ctx context.Context) error {
	_, err := s.db.ExecEx(ctx, "delete from otus.events where end_time < current_timestamp - interval '1 year';", nil)
	return err
}
