package tododomain

import "time"

type DueDate time.Time

func NewDueDate(date time.Time) (DueDate, error) {
	// TODO: validation

	return DueDate(date), nil
}

func (t DueDate) Value() time.Time {
	return time.Time(t)
}
