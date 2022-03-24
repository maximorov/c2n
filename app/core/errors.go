package core

var DefClientError = NewDefClientError()

func NewClientError(msg string) *ClientError {
	return &ClientError{msg}
}

func NewDefClientError() *ClientError {
	return &ClientError{`Команда не зрозуміла. Виберіть одну з тих, що нижче. ` + SymbLoopDown}
}

type ClientError struct {
	msg string
}

func (s *ClientError) Error() string {
	return s.msg
}
