package sqlstorage

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/jackc/pgx"
	"github.com/trb1maker/otus_golang_home_work/hw12_13_14_15_calendar/internal/storage"
)

var ErrEnvStorage = errors.New("server environment not set")

type Storage struct {
	db *pgx.Conn
}

func New() *Storage {
	return &Storage{}
}

func (s *Storage) initConfig() (*pgx.ConnConfig, error) {
	host, ok := os.LookupEnv("DBHOST")
	if !ok {
		return nil, ErrEnvStorage
	}

	port, ok := os.LookupEnv("DBPORT")
	if !ok {
		return nil, ErrEnvStorage
	}
	parsedPort, err := strconv.Atoi(port)
	if err != nil {
		return nil, err
	}

	dbname, ok := os.LookupEnv("DBNAME")
	if !ok {
		return nil, ErrEnvStorage
	}

	user, ok := os.LookupEnv("DBUSER")
	if !ok {
		return nil, ErrEnvStorage
	}

	password, ok := os.LookupEnv("DBPASSWORD")
	if !ok {
		return nil, ErrEnvStorage
	}

	config := &pgx.ConnConfig{
		Host:     host,
		Port:     uint16(parsedPort),
		Database: dbname,
		User:     user,
		Password: password,
	}

	return config, nil
}

func (s *Storage) Connect(ctx context.Context) error {
	config, err := s.initConfig()
	if err != nil {
		return err
	}

	s.db, err = pgx.Connect(*config)
	if err != nil {
		return err
	}

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
			from otus.events where event_id = $1 and is_deleted = false
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
