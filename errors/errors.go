package errors

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
