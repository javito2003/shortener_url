package app

type ErrorType int

const (
	ErrNotFound ErrorType = iota
	ErrConflict
	ErrInternal
	ErrUnauthorized
	ErrBadRequest
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

func NewUnauthorizedError(msg string) *AppError {
	return NewAppError(ErrUnauthorized, msg)
}

func NewConflictError(msg string) *AppError {
	return NewAppError(ErrConflict, msg)
}

func NewBadRequestError(msg string) *AppError {
	return NewAppError(ErrBadRequest, msg)
}
