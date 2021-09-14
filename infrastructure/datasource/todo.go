package datasource

import "time"

type Todo struct {
	ID                 int       `db:"id"`
	Title              string    `db:"title"`
	ImplementationDate time.Time `db:"implementation_date"`
	DueDate            time.Time `db:"due_date"`
	StatusID           uint      `db:"status_id"`
	PriorityID         uint      `db:"priority_id"`
	Memo               string    `db:"memo"`
}
