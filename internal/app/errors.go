package app

type ErrorType int

const (
	ErrNotFound ErrorType = iota
	ErrConflict
	ErrInternal
)

type AppError struct {
	Type    ErrorType
	Message string
}

func (e *AppError) Error() string {
	return e.Message
}

func NewAppError(t ErrorType, msg string) *AppError {
	return &AppError{
		Type:    t,
		Message: msg,
	}
}

func NewNotFoundError(msg string) *AppError {
	return NewAppError(ErrNotFound, msg)
}
