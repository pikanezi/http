package http

import (
	"github.com/gorilla/mux"
	"net/http"
)

// RouteHandler contains the main Handler and the given mux.Route.
// When an interceptor has been given, it fires them.
// If an Error occurs in a Before interceptor, it stops any further interceptor and writes the
// Error to the client.
// If an Error occurs in the Handler, it fires the OnError interceptors and the After interceptors.
// Any Error occurring in a OnError or After interceptor is ignored.
type RouteHandler struct {
	Route        *mux.Route
	handlerFunc  HandlerFunc
	interceptors []*interceptor
}

// finalHandler returns the handler to be used when applying any interceptor.
func (h *RouteHandler) finalHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w2, r2 := CreateResponseWriter(w), CreateRequest(r)
		if err := interceptors(h.interceptors).Before(w2, r2); err != nil {
			w2.WriteError(err)
			return
		}
		if err := h.handlerFunc(w2, r2); err != nil {
			w2.WriteError(err)
			interceptors(h.interceptors).OnError(w2, r2)
		}
		interceptors(h.interceptors).After(w2, r2)
	}
}

// Before add interceptors that must be run before the Request has been handled.
func (h *RouteHandler) Before(handlers ...HandlerFunc) *RouteHandler {
	for _, handler := range handlers {
		h.interceptors = append(h.interceptors, &interceptor{before, handler})
	}
	h.Route.Handler(h.finalHandler())
	return h
}

// After add interceptors that must be run after the Request has been handled.
func (h *RouteHandler) After(handlers ...HandlerFunc) *RouteHandler {
	for _, handler := range handlers {
		h.interceptors = append(h.interceptors, &interceptor{after, handler})
	}
	h.Route.Handler(h.finalHandler())
	return h
}

// OnError add interceptors that must be run when an Error occurs.
func (h *RouteHandler) OnError(handlers ...HandlerFunc) *RouteHandler {
	for _, handler := range handlers {
		h.interceptors = append(h.interceptors, &interceptor{onError, handler})
	}
	h.Route.Handler(h.finalHandler())
	return h
}

type when int

const (
	// before that the request is handled.
	before when = iota

	// after that the request has been handled.
	after

	//onError that the request returned an error.
	onError
)

type interceptor struct {
	when        when
	handlerFunc HandlerFunc
}

type interceptors []*interceptor

// Before runs every interceptors that must run before handling the Request.
// If an Error happens, it returns it as the request must be stop.
func (i interceptors) Before(w ResponseWriter, r *Request) *Error {
	for _, interceptor := range i {
		if interceptor.when == before {
			if err := interceptor.handlerFunc(w, r); err != nil {
				return err
			}
		}
	}
	return nil
}

// OnError runs every interceptors that must run when the Handler returns an Error.
func (i interceptors) OnError(w ResponseWriter, r *Request) {
	for _, interceptor := range i {
		if interceptor.when == onError {
			interceptor.handlerFunc(w, r)
		}
	}
}

// After runs every interceptors that must run when the Handler returns.
// If an Error in the Handler occurs, these interceptors are still fired.
func (i interceptors) After(w ResponseWriter, r *Request) {
	for _, interceptor := range i {
		if interceptor.when == after {
			interceptor.handlerFunc(w, r)
		}
	}
}

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
