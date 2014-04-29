package http

var (
	debugMode = false
)

func SetDebug(debug bool) {
	debugMode = debug
}

type Error struct {
	Error      string `json:"error,omitempty"`
	HttpCode   int    `json:"httpCode,omitempty"`
	StatusCode int    `json:"statusCode,omitempty"`
}

func NewError(err error, statusCode, httpCode int) *Error {
	return &Error{err.Error(), httpCode, statusCode}
}
