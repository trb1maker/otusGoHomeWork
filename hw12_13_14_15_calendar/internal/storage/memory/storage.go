package memorystorage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"                                                      //noling:gci
	"github.com/trb1maker/otus_golang_home_work/hw12_13_14_15_calendar/internal/storage" //noling:gci
)

const sqliteConnect = "file::memory:?cache=shared"

type Storage struct {
	db *sqlx.DB
}

func New() (*Storage, error) {
	db, err := sqlx.Connect("sqlite3", sqliteConnect)
	if err != nil {
		return nil, err
	}
	return &Storage{db: db}, db.Ping()
}

func (s *Storage) Connect(ctx context.Context) error {
	_, err := s.db.ExecContext(
		ctx,
		`
			create table users (
				user_id varchar primary key
			);

			create table events (
				event_id    varchar   primary key,
				title       varchar   not null,
				start_time  timestamp not null,
				end_time    timestamp not null,
				description varchar   not null,
				user_id     varchar   not null references users(user_id),
				notify      interval  not null,
				is_deleted  boolean   not null default false
			);

			create index events_user_id_idx on events (user_id);
		`,
	)
	return err
}

func (s *Storage) Close(_ context.Context) error {
	return s.db.Close()
}

func (s *Storage) InsertOne(ctx context.Context, event storage.Event) (eventID string, err error) {
	row := s.db.QueryRowContext(
		ctx,
		`
			insert into events (event_id, title, start_time, end_time, description, user_id, notify)
			values ($1, $2, $3, $4, $5, $6, $7)
			returning event_id
		`,
		uuid.New().String(),
		event.Title,
		event.StartTime.UTC(),
		event.EndTime.UTC(),
		event.Description,
		event.OwnerID,
		event.Notify,
	)
	return eventID, row.Scan(&eventID)
}

func (s *Storage) SelectOne(ctx context.Context, eventID string) (storage.Event, error) {
	result := s.db.QueryRowContext(
		ctx,
		`
			select title, start_time, end_time, description, user_id, notify
			from events
			where event_id = $1
			  and is_deleted = false
		`,
		eventID,
	)
	var e storage.Event
	if err := result.Scan(&e.Title, &e.StartTime, &e.EndTime, &e.Description, &e.OwnerID, &e.Notify); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return storage.Event{}, storage.ErrNotFound
		}
		return storage.Event{}, fmt.Errorf("can't scan result %w", err)
	}
	e.ID = eventID
	return e, nil
}

func (s *Storage) UpdateOne(ctx context.Context, event storage.Event) error {
	result, err := s.db.ExecContext(
		ctx,
		`
			update events set
				title = $1,
				start_time = $2,
				end_time = $3,
				description = $4,
				notify = $5
			where event_id = $6
		`,
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

	if cnt, _ := result.RowsAffected(); cnt == 0 {
		return storage.ErrNotFound
	}

	return nil
}

func (s *Storage) DeleteOne(ctx context.Context, eventID string) error {
	result, err := s.db.ExecContext(
		ctx,
		`
			update events
			set is_deleted = true
			where event_id = $1
		`,
		eventID,
	)
	if err != nil {
		return fmt.Errorf("can't delete event %w", err)
	}

	if cnt, _ := result.RowsAffected(); cnt == 0 {
		return storage.ErrNotFound
	}

	return nil
}

func (s *Storage) SelectAllEvents(ctx context.Context, userID string) ([]storage.Event, error) {
	rows, err := s.db.QueryContext(
		ctx,
		`
			select event_id, title, start_time, end_time, description, notify
			from events
			where is_deleted = false 
			  and user_id = $1
		`,
		userID,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
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
	result := s.db.QueryRowContext(
		ctx,
		`
			select event_id, title, start_time, end_time, description, notify
			from events
			where is_deleted = false
			  and user_id = $1
			  and start_time >= current_timestamp
			order by start_time
			limit 1
		`,
		userID,
	)

	var (
		eventID, title, description string
		startTime, endTime          time.Time
		notify                      time.Duration
	)
	if err := result.Scan(&eventID, &title, &startTime, &endTime, &description, &notify); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
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
	ownerID string,
	from time.Time,
	to time.Time,
) ([]storage.Event, error) {
	rows, err := s.db.QueryContext(
		ctx,
		`
			select event_id, title, start_time, end_time, description, notify
			from events
			where is_deleted = false
			  and user_id = $1
			  and start_time >= $2
			  and start_time <= $3
		`,
		ownerID,
		from,
		to,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
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
			OwnerID:     ownerID,
		})
	}

	return events, nil
}

func (s *Storage) RegisterUser(ctx context.Context) (userID string, err error) {
	userID = uuid.New().String()
	if _, err := s.db.ExecContext(ctx, "insert into users (user_id) values ($1)", userID); err != nil {
		return "", err
	}
	return userID, err
}
