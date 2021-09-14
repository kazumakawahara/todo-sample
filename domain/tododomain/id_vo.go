package tododomain

type ID int

func NewID(id int) (ID, error) {
	// TODO: validation

	return ID(id), nil
}

func (i ID) Value() int {
	return int(i)
}
