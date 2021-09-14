package tododomain

type Memo string

func NewMemo(memo string) (Memo, error) {
	// TODO: validation

	return Memo(memo), nil
}

func (p Memo) Value() string {
	return string(p)
}
