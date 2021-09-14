package tododomain

type Status uint

const (
	TODO Status = iota + 1
	DOING
	DONE
)

func NewStatus(statusID uint) (Status, error) {
	// TODO: validation

	return Status(statusID), nil
}

func (s Status) Value() uint {
	return uint(s)
}
