package tododomain

import "time"

type ImplementationDate time.Time

func NewImplementationDate(date time.Time) (ImplementationDate, error) {
	// TODO: validation

	return ImplementationDate(date), nil
}

func (t ImplementationDate) Value() time.Time {
	return time.Time(t)
}
