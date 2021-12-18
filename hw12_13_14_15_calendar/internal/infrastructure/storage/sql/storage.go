package sqlstorage

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/alexandr-lakeev/otus-home-work/hw12_13_14_15_calendar/internal/domain"
	"github.com/alexandr-lakeev/otus-home-work/hw12_13_14_15_calendar/internal/domain/models"
	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
)

type Storage struct {
	db   *sqlx.DB
	conn *sql.Conn
	dsn  string
}

type dbEvent struct {
	ID          uuid.UUID     `db:"id"`
	UserID      uuid.UUID     `db:"user_id"`
	Date        time.Time     `db:"date"`
	Duration    time.Duration `db:"duration"`
	Title       string        `db:"title"`
	Description string        `db:"description"`
}

func New(dsn string) *Storage {
	return &Storage{
		dsn: dsn,
	}
}

func (s *Storage) Connect(ctx context.Context) (err error) {
	s.db, err = sqlx.Open("pgx", s.dsn)
	if err != nil {
		return err
	}

	s.conn, err = s.db.Conn(ctx)
	if err != nil {
		return nil
	}

	return s.db.PingContext(ctx)
}

func (s *Storage) Close(ctx context.Context) error {
	return s.conn.Close()
}

func (s *Storage) Get(ctx context.Context, id uuid.UUID) (*models.Event, error) {
	query := `
		SELECT 
			id, user_id, title, date, duration, description
		FROM 
			events 
		WHERE 
			id = $1
	`

	rows, err := s.db.QueryxContext(ctx, query, id)
	if err != nil {
		return nil, err
	}

	if !rows.Next() {
		return nil, domain.ErrEventNotFound
	}

	var event dbEvent

	err = rows.StructScan(&event)
	if err != nil {
		return nil, err
	}

	err = rows.Close()
	if err != nil {
		return nil, err
	}

	domainEvent, err := dbToDomainEvent(event)
	if err != nil {
		return nil, err
	}

	return &domainEvent, nil
}

func (s *Storage) Add(ctx context.Context, event *models.Event) error {
	query := `
		INSERT INTO events (id, user_id, title, date, duration, description)
		VALUES (:id, :userID, :title, :date, :duration, :description)
	`

	_, err := s.db.NamedExecContext(ctx, query, map[string]interface{}{
		"id":          event.ID.String(),
		"UserID":      event.UserID.String(),
		"title":       event.Title,
		"date":        event.Date,
		"duration":    event.Duration.Minutes(),
		"description": event.Description,
	})

	return err
}

func (s *Storage) Update(ctx context.Context, event *models.Event) error {
	query := `
		UPDATE events
		SET
			title = :title,
			date = :date,
			duration = :duration,
			description = :description,
			updated_at = NOW()
		WHERE
			id = id AND
			user_id = :userID
	`

	_, err := s.db.NamedExecContext(ctx, query, map[string]interface{}{
		"id":          event.ID.String(),
		"userID":      event.UserID.String(),
		"title":       event.Title,
		"date":        event.Date,
		"duration":    event.Duration.Minutes(),
		"description": event.Description,
	})

	return err
}

func (s *Storage) GetList(ctx context.Context, userID models.ID, from, to time.Time) ([]models.Event, error) {
	query := `
		SELECT 
			id, user_id, title, date, duration, description
		FROM 
			events 
		WHERE 
			user_id = :userID
			AND date BETWEEN :from AND :to
		ORDER BY date
	`

	rows, err := s.db.NamedQueryContext(ctx, query, map[string]interface{}{
		"userID": userID,
		"from":   from.Format("2006-01-02"),
		"to":     to.Format("2006-01-02"),
	})

	if err != nil {
		return nil, err
	}

	if !rows.Next() {
		return nil, domain.ErrEventNotFound
	}

	var events []models.Event

	for rows.Next() {
		var event dbEvent

		err := rows.StructScan(&event)
		if err != nil {
			return nil, err
		}

		domainEvent, err := dbToDomainEvent(event)
		if err != nil {
			return nil, err
		}

		events = append(events, domainEvent)
	}

	err = rows.Close()
	if err != nil {
		return nil, err
	}

	return events, nil
}

func dbToDomainEvent(event dbEvent) (models.Event, error) {
	duration, err := time.ParseDuration(fmt.Sprintf("%dm", event.Duration))
	if err != nil {
		return models.Event{}, err
	}

	return models.Event{
		ID:          event.ID,
		UserID:      event.UserID,
		Title:       event.Title,
		Date:        event.Date,
		Duration:    duration,
		Description: event.Description,
	}, nil
}
