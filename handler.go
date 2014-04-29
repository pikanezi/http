package http

import "net/http"

type HandlerFunc func(ResponseWriter, *Request) *Error

func createHandler(handler HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w2, r2 := createResponseWriter(w), createRequest(r)
		if err := handler(w2, r2); err != nil {
			w2.WriteError(err)
		}
	}
}

func Handle(route string, handler http.Handler) {
	http.Handle(route, handler)
}

func ListenAndServe(addr string, handler http.Handler) error {
	return http.ListenAndServe(addr, handler)
}
