/*
Package http overrides net/http and gorilla/pat golang library to use JSON with Request and ResponseWriter in a more easy way.

Its purpose is to create web server sending and receiving JSON data in few lines.

Full Example

	import (
    	"github.com/pikanezi/http"
    	"log"
	)

	type Object struct {
    	SomeField string `json:"someField,omitempty"`
	}

	func HelloWorldHandler(w http.ResponseWriter, r *http.Request) *http.Error {
    	object := &Object{"Hello World"}
    	if err := w.WriteJSON(object); err != nil {
        	return &http.Error{Error: err.Error()
    	}
	}

	func main() {

    	// NewRouter takes the domain to authorize it cross-domains requests
    	r := http.NewRouter()

    	r.Get("/hello/world", HelloWorldHandler)

    	log.Fatal(http.ListenAndServe(":8080", r)
	}

*/
package http

var (
	debugMode = false
)

// SetDebug set the debug at true so it prints some debug.
func SetDebug(debug bool) {
	debugMode = debug
}

// Error is a type that must be used when returning from the Handlers.
type Error struct {
	Error      string `json:"error,omitempty"`
	HTTPCode   int    `json:"httpCode,omitempty"`
	StatusCode int    `json:"statusCode,omitempty"`
}

// NewErrorAPI returns a new Error with a statusCode.
func NewErrorAPI(err error, statusCode, httpCode int) *Error {
	return &Error{err.Error(), httpCode, statusCode}
}

// NewError returns a new Error.
func NewError(err error, httpCode int) *Error {
	return &Error{
		Error:    err.Error(),
		HTTPCode: httpCode,
	}
}
