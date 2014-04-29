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
	hooks []HandlerFunc
}

func (self *Router) Domain() string          { return self.domain }
func (self *Router) SetDomain(domain string) { self.domain = domain }

// Returns a new router with the given path
func NewRouter(domain string) *Router {
	return &Router{pat.New(), domain, make([]HandlerFunc, 0)}
}

// Add a function to be executed before serving HTTP
func (self *Router) AddHooks(hooks ...HandlerFunc) {
	self.hooks = append(self.hooks, hooks...)
}

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
	wr := createResponseWriter(w)
	rr := createRequest(r)
	if strings.ToLower(r.Method) == "options" {
		http.Redirect(wr, r, r.RequestURI, 200)
		return
	}
	if err := self.runHooks(wr, rr); err != nil {
		wr.WriteError(err)
		return
	}
	self.Router.ServeHTTP(w, r)
}
