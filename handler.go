package http

import "net/http"

// HandlerFunc must returns an Error that will be handled by the Router itself.
// See implementation of the ServeHTTP method of the Router.
type HandlerFunc func(ResponseWriter, *Request) *Error

func createHandler(handler HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w2, r2 := CreateResponseWriter(w), CreateRequest(r)
		if err := handler(w2, r2); err != nil {
			w2.WriteError(err)
		}
	}
}

// Handle registers the handler for the given pattern in the golang http DefaultServeMux.
// This should be avoided, use ListenAndServe instead.
func Handle(route string, handler http.Handler) {
	http.Handle(route, handler)
}

// ListenAndServe listens on the TCP network address addr and then calls Serve with handler to handle requests on incoming connections.
// Handler is typically the Router itself.
func ListenAndServe(addr string, handler http.Handler) error {
	return http.ListenAndServe(addr, handler)
}
