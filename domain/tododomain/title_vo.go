package tododomain

type Title string

func NewTitle(title string) (Title, error) {
	// TODO: validation

	return Title(title), nil
}

func (t Title) Value() string {
	return string(t)
}
