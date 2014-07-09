package http

import (
	"fmt"
	"io"
	"io/ioutil"
	"testing"
)

func ExampleNewRouter() {
	router := NewRouter()

	router.Get("/", func(w ResponseWriter, r *Request) *Error {
		w.WriteJSON("Hello!")
		return nil
	})

	router.Get("/admin", func(w ResponseWriter, r *Request) *Error {
		// do stuff
		return nil
	}).Before(func(w ResponseWriter, r *Request) *Error {
		// check if user is an admin
		return nil
	})

	router.Post("/upload", func(w ResponseWriter, r *Request) *Error {
		if err := r.ForEachFile("file", func(index int, reader io.Reader) error {
			// Do something with the reader
			return nil
		}); err != nil {
			return NewError(err, 400)
		}
		return nil
	})
	ListenAndServe(":8080", router)
}

func TestHandler(t *testing.T) {
	mainHandler := HandlerFunc(func(w ResponseWriter, r *Request) *Error {
		fmt.Println("mainHandler")
		w.WriteJSON("mainHandler")
		return nil
	})
	beforeMainHandler := HandlerFunc(func(w ResponseWriter, r *Request) *Error {
		fmt.Println("before Handler!")
		return nil
	})
	afterMainHandler := HandlerFunc(func(w ResponseWriter, r *Request) *Error {
		fmt.Println("after Handler!")
		return nil
	})

	filesHander := HandlerFunc(func(w ResponseWriter, r *Request) *Error {
		if err := r.ForEachFile("file", func(index int, reader io.Reader) error {
			buffer, err := ioutil.ReadAll(reader)
			if err != nil {
				return err
			}
			fmt.Println(buffer)
			return nil
		}); err != nil {
			return NewError(err, 400)
		}
		w.WriteJSON("OK")
		return nil
	})

	r := NewRouter()
	r.Get("/", mainHandler).Before(beforeMainHandler).After(afterMainHandler)
	r.Post("/files", filesHander)

	ListenAndServe(":5555", r)
}
