package input

import "time"

type Todo struct {
	ID                 int       `json:"id"`
	Title              string    `json:"title"`
	ImplementationDate time.Time `json:"implementationDate"`
	DueDate            time.Time `json:"dueDate"`
	StatusID           uint      `json:"statusID"`
	PriorityID         uint      `json:"priorityID"`
	Memo               string    `json:"memo"`
}
