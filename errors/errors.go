package errors

type AppError struct {
	Code    int
	Message string
	Err     error
}

func NewAppError(code int, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

func (e *AppError) Error() string {
	return e.Message
}

type UserNotFoundError struct{}

func (u UserNotFoundError) Error() string {
	return "user not found"
}

type UserAlreadyExists struct{}

func (u UserAlreadyExists) Error() string {
	return "user already exists"
}

type MissingConfigError struct {
	Message string
}

func (m MissingConfigError) Error() string {
	return m.Message
}
