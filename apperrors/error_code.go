package apperrors

type code string

const (
	InvalidParameterCode    code = "InvalidParameter"
	InternalServerErrorCode code = "InternalServerError"
	TodoNotFoundCode        code = "TodoNotFound"
)

func (c code) value() string {
	return string(c)
}
