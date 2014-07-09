package http

import (
	"testing"
	"fmt"
	"github.com/pikanezi/http"
)

func ExampleNewRouter() {
	router := http.NewRouter()

	router.Get("/", func(w http.ResponseWriter, r *http.Request) *http.Error {
			w.WriteJSON("Hello!")
			return nil
		})

	router.Get("/admin", func(w http.ResponseWriter, r *http.Request) *http.Error {
			// check user is admin
			return nil
		})
	http.ListenAndServe(":8080", router)
}

func TestHandler(t *testing.T) {
	mainHandler := HandlerFunc(func(w ResponseWriter, r *Request) *Error {
		fmt.Println("mainHandler")
		w.WriteJSON("mainHandler")
		return nil
	})
	beforeMainHandler := HandlerFunc(func(w ResponseWriter, r *Request) *Error{
		fmt.Println("before Handler!")
		return nil
	})
	afterMainHandler := HandlerFunc(func(w ResponseWriter, r *Request) *Error{
		fmt.Println("after Handler!")
		return nil
	})
	r := NewRouter()
	r.Get("/", mainHandler).Before(beforeMainHandler).After(afterMainHandler)
	ListenAndServe(":5555", r)
}
