package http

import (
	"github.com/gorilla/pat"
	"net/http"
	"strings"
)

// Router embeds pat.Router.
// Router is a request router that implements a pat-like API.
// pat docs: http://godoc.org/github.com/bmizerany/pat
type Router struct {
	*pat.Router

	customHeaders Header
}

// Header represents custom header to be set to the response before .
type Header map[string]string

// Add add a key-value pair to the header.
func (h Header) Add(key, value string) { h[key] = value }

// SetCustomHeader set the customHeader of the router.
func (router *Router) SetCustomHeader(customHeader Header) { router.customHeaders = customHeader }

// CustomHeader returns the customHeaders of the router.
func (router *Router) CustomHeader() Header { return router.customHeaders }

// NewRouter returns a new router with the given domain.
func NewRouter() *Router {
	return &Router{
		Router:        pat.New(),
		customHeaders: nil,
	}
}

func (router *Router) add(method, route string, h HandlerFunc) *RouteHandler {
	handler := &RouteHandler{
		Route:        router.Add(method, route, createHandler(h)),
		handlerFunc:  h,
		interceptors: make([]*interceptor, 0),
	}
	return handler
}

// Get registers a pattern with a handler for GET requests.
func (router *Router) Get(route string, h HandlerFunc) *RouteHandler {
	return router.add("GET", route, h)
}

// Post registers a pattern with a handler for POST requests.
func (router *Router) Post(route string, h HandlerFunc) *RouteHandler {
	return router.add("POST", route, h)
}

// Delete registers a pattern with a handler for DELETE requests.
func (router *Router) Delete(route string, h HandlerFunc) *RouteHandler {
	return router.add("DELETE", route, h)
}

// Put registers a pattern with a handler for PUT requests.
func (router *Router) Put(route string, h HandlerFunc) *RouteHandler {
	return router.add("PUT", route, h)
}

// ServeHTTP dispatches the handler registered in the matched route.
// It performs any hooks and add the domain registered in the Router to be allowed for cross-domain requests.
func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	wr := CreateResponseWriter(w)
	wr.addCustomPreHeader(router.customHeaders)
	if strings.ToLower(r.Method) == "options" {
		http.Redirect(wr, r, r.RequestURI, 200)
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
