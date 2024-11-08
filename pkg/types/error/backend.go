package yError

type NoBackendError struct {
	msg string
}

func (e *NoBackendError) Error() string {
	return e.msg
}

func NewNoBackendError(msg string) *NoBackendError {
	return &NoBackendError{msg: msg}
}
