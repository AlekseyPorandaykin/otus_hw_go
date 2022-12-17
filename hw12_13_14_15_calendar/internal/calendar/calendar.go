package calendar

import (
	"time"
)

type Event struct {
	ID            string    `db:"id"`
	Title         string    `db:"title"`
	DateTimeStart time.Time `db:"datetime_from"`
	DateTimeEnd   time.Time `db:"datetime_to"`
	Description   string    `db:"description"`
	CreatedBy     UserID    `db:"created_by"`
	RemindFrom    time.Time `db:"start_notify"`
}

type UserID int

type Notification struct {
	ID       string
	Title    string
	Datetime time.Time
	UserTo   UserID
}
