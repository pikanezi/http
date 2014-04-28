package main

import (
	"github.com/gorilla/mux"
	"github.com/gorilla/pat"
	"net/http"
	"strings"
)

type Router struct {
	*pat.Router
	domain string
}

func (self *Router) Domain() string          { return self.domain }
func (self *Router) SetDomain(domain string) { self.domain = domain }

// Returns a new router with the given path
func NewRouter(domain string) *Router {
	return &Router{pat.New(), domain}
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
	if strings.ToLower(r.Method) == "options" {
		http.Redirect(wr, r, r.RequestURI, 200)
		return
	}
	self.Router.ServeHTTP(w, r)
}
