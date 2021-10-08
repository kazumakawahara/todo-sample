package tododomain

import (
	"github.com/kazumakawahara/todo-sample/apperrors"
	"unicode/utf8"
)

type Title string

func NewTitle(title string) (Title, error) {

	// TODO: validation
	if utf8.RuneCountInString(title) > 10 {
		return "", apperrors.InvalidParameter
	}

	return Title(title), nil
}

func (t Title) Value() string {
	return string(t)
}
