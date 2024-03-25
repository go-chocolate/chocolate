package errorx

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *Error) Error() string {
	return e.Message
}

type TextError string

func (e TextError) Error() string {
	return string(e)
}

func Code(code int, msg ...string) *Error {
	e := &Error{Code: code}
	if len(msg) > 0 {
		e.Message = msg[0]
	}
	return e
}
