/*
Override of net/http and gorilla/pat golang library to use JSON with Request and ResponseWritter in a more easy way.

Its purpose is to create web server sending and receiving JSON data in few lines.

Full Example

	import (
    	"github.com/pikanezi/http"
    	"log"
	)

	const (
    	KEY = "SECRET_KEY"
	)

	func SecureHook(w http.ResponseWriter, r *http.Request) *http.Error {
    	if r.Header.Get("x-api-key") != KEY {
        	return &http.Error{Error:"Wrong API Key", HttpCode: 403}
    	}
    	return nil
	}

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
    	r := http.NewRouter("example.com")

    	// Add a Hook, every Hooks will be executed before executing an Handler
    	r.AddHooks(SecureHook)

    	r.Get("/hello/world", HelloWorldHandler)

    	log.Fatal(http.ListenAndServe(":8080", r)
	}

*/
package http

var (
	debugMode = false
)

func SetDebug(debug bool) {
	debugMode = debug
}

// The error type that must be used when returning from the Handlers.
type Error struct {
	Error      string `json:"error,omitempty"`
	HttpCode   int    `json:"httpCode,omitempty"`
	StatusCode int    `json:"statusCode,omitempty"`
}

// Create a new Error.
func NewError(err error, statusCode, httpCode int) *Error {
	return &Error{err.Error(), httpCode, statusCode}
}
