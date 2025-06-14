package apiError

type Error interface {
	Message() string
	Code() int
	Error() string
}

type err struct {
	msg  string
	code int
}

func (e err) Error() string {
	return e.Message()
}

func (e err) Message() string {
	return e.msg
}

func (e err) Code() int {
	return e.code
}

func New(code int, msg string) Error {
	return err{
		msg:  msg,
		code: code,
	}
}
