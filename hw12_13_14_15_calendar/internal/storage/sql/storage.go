package sqlstorage

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/AlekseyPorandaykin/otus_hw_go/hw12_13_14_15_calendar/internal/calendar"
	"github.com/AlekseyPorandaykin/otus_hw_go/hw12_13_14_15_calendar/internal/config"
	"github.com/jackc/pgx"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
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

func (s *SQLStorage) Connect(ctx context.Context) error {
	var err error
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

	return err
}

const createQuery = `
INSERT INTO 
    public.events(id,title,description,datetime_from,datetime_to,created_by,start_notify,notify_status,created_at)
	VALUES 
	    (:id,:title,:description,:datetime_from,:datetime_to,:created_by,:start_notify,:notify_status,:created_at)
RETURNING id
	`

func (s *SQLStorage) CreateEvent(ctx context.Context, e *calendar.Event) error {
	var id string
	smtp, err := s.db.PrepareNamed(createQuery)
	if err != nil {
		return err
	}
	defer smtp.Close()

	return smtp.GetContext(ctx, &id, e)
}

const updateQuery = `
	UPDATE public.events 
		SET title = :title, 
			description = :description,
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

func (s *SQLStorage) EventByID(ctx context.Context, id string) (*calendar.Event, error) {
	query := `SELECT * FROM public.events WHERE id=$1`
	row := s.db.QueryRowxContext(ctx, query, id)
	if row.Err() != nil {
		return nil, row.Err()
	}
	var e calendar.Event
	err := row.StructScan(&e)

	if err != nil && strings.Contains(err.Error(), pgx.ErrNoRows.Error()) {
		return nil, nil
	}
	return &e, err
}

func (s *SQLStorage) EventsByPeriod(ctx context.Context, start, end time.Time, limit int) ([]*calendar.Event, error) {
	var (
		events []*calendar.Event
		query  = `SELECT * FROM public.events WHERE datetime_from>=$1 AND datetime_to <= $2 LIMIT $3`
	)
	rows, err := s.db.QueryxContext(ctx, query, start, end, limit)
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

func (s *SQLStorage) GetEventsWithRemindStatus(
	ctx context.Context, time time.Time, status calendar.RemindStatus,
) ([]*calendar.Event, error) {
	var (
		events []*calendar.Event
		query  = `SELECT * FROM public.events WHERE start_notify<=$1 AND notify_status = $2`
	)
	rows, err := s.db.QueryxContext(ctx, query, time, status)
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

func (s *SQLStorage) UpdateRemindStatus(ctx context.Context, id string, status calendar.RemindStatus) error {
	query := `UPDATE public.events SET notify_status=$1 WHERE id=$2`
	res, err := s.db.ExecContext(ctx, query, status, id)
	if err != nil {
		return err
	}
	rowA, err := res.RowsAffected()
	if rowA == 0 {
		return ErrNotUpdate
	}
	return err
}

func (s *SQLStorage) GetOldEventIDs(ctx context.Context, oldTime time.Time) ([]string, error) {
	var (
		query = `SELECT id FROM public.events WHERE datetime_to <= $1`
		ids   = make([]string, 0)
	)
	err := s.db.SelectContext(ctx, &ids, query, oldTime)

	return ids, err
}

func (s *SQLStorage) DeleteEventByIDs(ctx context.Context, ids []string) error {
	query := `DELETE FROM public.events WHERE id = ANY($1)`
	_, err := s.db.ExecContext(ctx, query, pq.Array(ids))
	return err
}

const saveLogQuery = `
INSERT INTO 
    public.logs(event_id, body, created_at)
	VALUES 
	    ($1,$2, $3)
	`

func (s *SQLStorage) Save(ctx context.Context, eventID string, body []byte, date time.Time) error {
	_, err := s.db.ExecContext(ctx, saveLogQuery, eventID, body, date)
	return err
}

func (s *SQLStorage) Close(ctx context.Context) error {
	if s.db == nil {
		return ErrNotInitDB
	}
	return s.db.Close()
}
