package http

var (
	debugMode = false
)

func setDebug(debug bool) {
	debugMode = debug
}

type Error struct {
	Error      string `json:"error,omitempty"`
	StatusCode int    `json:"statusCode,omitempty"`
	HttpCode   int    `json:"httpCode,omitempty"`
}

func NewError(err error, statusCode, httpCode int) *Error {
	return &Error{err.Error(), statusCode, httpCode}
}
