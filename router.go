package http

import (
	"github.com/gorilla/mux"
	"github.com/gorilla/pat"
	"net/http"
	"strings"
)

type Router struct {
	*pat.Router
	domain string
	hooks  []HandlerFunc
}

func (self *Router) Domain() string {
	return self.domain
}

func (self *Router) SetDomain(domain string) {
	self.domain = domain
}

// Returns a new router with the given domain.
func NewRouter(domain string) *Router { return &Router{pat.New(), domain, make([]HandlerFunc, 0)} }

// Add a function to be executed before serving HTTP.
func (self *Router) AddHooks(hooks ...HandlerFunc) { self.hooks = append(self.hooks, hooks...) }

func (self *Router) runHooks(w ResponseWriter, r *Request) *Error {
	for _, hook := range self.hooks {
		if err := hook(w, r); err != nil {
			return err
		}
	}
	return nil
}

func (self *Router) Get(route string, h HandlerFunc) *mux.Route {
	return self.Add("GET", route, createHandler(h))
}

func (self *Router) Post(route string, h HandlerFunc) *mux.Route {
	return self.Add("POST", route, createHandler(h))
}

func (self *Router) Delete(route string, h HandlerFunc) *mux.Route {
	return self.Add("DELETE", route, createHandler(h))
}

func (self *Router) Put(route string, h HandlerFunc) *mux.Route {
	return self.Add("PUT", route, createHandler(h))
}

func (self *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	wr, rr := createResponseWriter(w), createRequest(r)
	if strings.ToLower(r.Method) == "options" {
		http.Redirect(wr, r, r.RequestURI, 200)
		return
	}
	wr.addCORSHeaders(self.domain)
	if err := self.runHooks(wr, rr); err != nil {
		wr.WriteError(err)
		return
	}
	self.Router.ServeHTTP(w, r)
}

// Redirect replies to the request with a redirect to url, which may be a path relative to the request path.
func Redirect(w ResponseWriter, r *Request, urlStr string, code int) {
	http.Redirect(w, r.Request, urlStr, code)
}

// StatusText returns a text for the HTTP status code. It returns the empty string if the code is unknown.
func StatusText(code int) string {
	return http.StatusText(code)
}
