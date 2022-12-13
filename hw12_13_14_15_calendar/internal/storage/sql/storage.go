package sqlstorage

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/AlekseyPorandaykin/otus_hw_go/hw12_13_14_15_calendar/internal/calendar"
	"github.com/AlekseyPorandaykin/otus_hw_go/hw12_13_14_15_calendar/internal/config"
	"github.com/jackc/pgx"
	"github.com/jmoiron/sqlx"
)

var (
	ErrNotInitDB = errors.New("not init db")
	ErrNotUpdate = errors.New("not update record")
)

type SQLStorage struct {
	conn *config.StorageConfig
	db   *sqlx.DB
}

func New(conn *config.StorageConfig) *SQLStorage {
	return &SQLStorage{
		conn: conn,
	}
}

func (s *SQLStorage) Connect(ctx context.Context) (err error) {
	s.db, err = sqlx.ConnectContext(
		ctx,
		"pgx",
		fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s",
			s.conn.User,
			s.conn.Password,
			s.conn.Host,
			s.conn.Port,
			s.conn.DBName,
		),
	)
	if err == nil {
		err = s.db.Ping()
	}

	return
}

const createQuery = `
	INSERT INTO public.events (id, title, datetime_from, datetime_to,created_by,start_notify)
		VALUES (:id, :title, :datetime_from, :datetime_to, :created_by, :start_notify)
	RETURNING id
	`

func (s *SQLStorage) CreateEvent(ctx context.Context, e *calendar.Event) (err error) {
	var id string
	smtp, err := s.db.PrepareNamed(createQuery)
	if err != nil {
		return
	}
	defer smtp.Close()
	err = smtp.GetContext(ctx, &id, e)

	return
}

const updateQuery = `
	UPDATE public.events 
		SET title = :title, 
		    datetime_from = :datetime_from, 
		    datetime_to = :datetime_to, 
		    created_by = :created_by, 
		    start_notify = :start_notify
	WHERE ID = :id;`

func (s *SQLStorage) UpdateEvent(ctx context.Context, e *calendar.Event) error {
	res, err := s.db.NamedExecContext(ctx, updateQuery, e)
	if err != nil {
		return err
	}
	rowA, err := res.RowsAffected()
	if rowA == 0 {
		return ErrNotUpdate
	}
	return err
}

const deleteEvent = `DELETE FROM public.events WHERE id=$1`

func (s *SQLStorage) DeleteEvent(ctx context.Context, id string) error {
	_, err := s.db.ExecContext(ctx, deleteEvent, id)
	return err
}

const getEventByID = `SELECT * FROM public.events WHERE id=$1`

func (s *SQLStorage) EventByID(ctx context.Context, id string) (*calendar.Event, error) {
	row := s.db.QueryRowxContext(ctx, getEventByID, id)
	if row.Err() != nil {
		return nil, row.Err()
	}
	var e calendar.Event
	err := row.StructScan(&e)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	return &e, err
}

const getEventsByPeriod = `SELECT * FROM public.events WHERE datetime_from>=$1 AND datetime_to <= $2 LIMIT $3`

func (s *SQLStorage) EventsByPeriod(ctx context.Context, start, end time.Time, limit int) ([]*calendar.Event, error) {
	var events []*calendar.Event
	rows, err := s.db.QueryxContext(ctx, getEventsByPeriod, start, end, limit)
	if err != nil {
		return events, err
	}
	defer rows.Close()
	for rows.Next() {
		e := calendar.Event{}
		if err := rows.StructScan(&e); err != nil {
			return events, err
		}
		events = append(events, &e)
	}

	return events, nil
}

func (s *SQLStorage) Close(ctx context.Context) error {
	if s.db == nil {
		return ErrNotInitDB
	}
	return s.db.Close()
}
