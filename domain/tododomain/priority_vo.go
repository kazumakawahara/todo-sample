package tododomain

type Priority uint

const (
	UNKNOWN Priority = iota + 1
	LOW
	MEDIUM
	HIGH
)

func NewPriority(priorityID uint) (Priority, error) {
	// TODO: validation

	return Priority(priorityID), nil
}

func (p Priority) Value() uint {
	return uint(p)
}
