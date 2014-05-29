package http

import (
	"github.com/gorilla/mux"
	"github.com/gorilla/pat"
	"net/http"
	"strings"
)

// Router embeds pat.Router.
// Router is a request router that implements a pat-like API.
// pat docs: http://godoc.org/github.com/bmizerany/pat
type Router struct {
	*pat.Router
	domain string
	hooks  []HandlerFunc
}

// Domain returns the domain of the Router.
func (router *Router) Domain() string {
	return router.domain
}

// SetDomain set the domain of the Router.
func (router *Router) SetDomain(domain string) {
	router.domain = domain
}

// NewRouter returns a new router with the given domain.
func NewRouter(domain string) *Router { return &Router{pat.New(), domain, make([]HandlerFunc, 0)} }

// AddHooks add a function to be executed before serving HTTP.
func (router *Router) AddHooks(hooks ...HandlerFunc) { router.hooks = append(router.hooks, hooks...) }

func (router *Router) runHooks(w ResponseWriter, r *Request) *Error {
	for _, hook := range router.hooks {
		if err := hook(w, r); err != nil {
			return err
		}
	}
	return nil
}

// Get registers a pattern with a handler for GET requests.
func (router *Router) Get(route string, h HandlerFunc) *mux.Route {
	return router.Add("GET", route, createHandler(h))
}

// Post registers a pattern with a handler for POST requests.
func (router *Router) Post(route string, h HandlerFunc) *mux.Route {
	return router.Add("POST", route, createHandler(h))
}

// Delete registers a pattern with a handler for DELETE requests.
func (router *Router) Delete(route string, h HandlerFunc) *mux.Route {
	return router.Add("DELETE", route, createHandler(h))
}

// Put registers a pattern with a handler for PUT requests.
func (router *Router) Put(route string, h HandlerFunc) *mux.Route {
	return router.Add("PUT", route, createHandler(h))
}

// ServeHTTP dispatches the handler registered in the matched route.
// It performs any hooks and add the domain registered in the Router to be allowed for cross-domain requests.
func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	wr, rr := createResponseWriter(w), createRequest(r)
	wr.addCORSHeaders(router.domain)
	if strings.ToLower(r.Method) == "options" {
		http.Redirect(wr, r, r.RequestURI, 200)
		return
	}
	if err := router.runHooks(wr, rr); err != nil {
		wr.WriteError(err)
		return
	}
	router.Router.ServeHTTP(w, r)
}

// Redirect replies to the request with a redirect to url, which may be a path relative to the request path.
func Redirect(w ResponseWriter, r *Request, urlStr string, code int) {
	http.Redirect(w, r.Request, urlStr, code)
}

// StatusText returns a text for the HTTP status code. It returns the empty string if the code is unknown.
func StatusText(code int) string {
	return http.StatusText(code)
}
