package http

import (
	"github.com/gorilla/mux"
	"github.com/gorilla/pat"
	"net/http"
	"strings"
	"fmt"
)

var (
	routers = []*Router{}
)

// ObjectRester
type ObjectRester interface {
	GetRootURL() string
	Get(ResponseWriter, *Request) *Error
	Post(ResponseWriter, *Request) *Error
	Put(ResponseWriter, *Request) *Error
	Delete(ResponseWriter, *Request) *Error
}

// Router is a pat.Router which has a domain for handling CORS requests
// and hooks to execute some functions before executing the HandlerFunc.
type Router struct {
	*pat.Router
	domain string
	hooks  []HandlerFunc
}

// NewRouter returns a new Router.
func NewRouter(domain string) *Router {
	router := &Router{pat.New(), domain, make([]HandlerFunc, 0)}
	routers = append(routers, router)
	return router
}

// Domain returns the domain of the Router.
func (self *Router) Domain() string {
	return self.domain
}

// SetDomain set the domain of the Router.
func (self *Router) SetDomain(domain string) {
	self.domain = domain
}

// Add a function to be executed before serving HTTP.
func (self *Router) AddHooks(hooks ...HandlerFunc) { self.hooks = append(self.hooks, hooks...) }

// Register add a POST, GET, PUT and DELETE HandlerFunc to add, get, edit and delete an instance of this object.
// Example:
//
// 		type User struct {
//			Name string `json:"name"`
//		}
//
// 		var (
//			Jonathan = &User{"Jonathan"}
//			Vincent = &User{"Vincent"}
//		)
//
//		type UserHandler struct {}
//
//		func (handler *UserHandler) GetRootURL() string { return "users" }
//
//		func (handler *UserHandler) Get(w http.ResponseWriter, r *http.Request) *http.Error {
//			// Do something
//		}
//
//		func (handler *UserHandler) Post(w http.ResponseWriter, r *http.Request) *http.Error {
//			// Do something
//		}
//
//		func (handler *UserHandler) Put(w http.ResponseWriter, r *http.Request) *http.Error {
//			// Do something
//		}
//
//		func (handler *UserHandler) Delete(w http.ResponseWriter, r *http.Request) *http.Error {
//			// Do something
//		}
//
//		func main() {
//			r := http.NewRouter("example.com")
//			r.Register(&UserHandler{})
//			http.ListenAndServe(":8080", r)
//		}
//
func (self *Router) Register(object ObjectRester) {
	route := fmt.Sprintf("/%v", object.GetRootURL())
	self.Get(route, object.Get)
	self.Post(route, object.Post)
	self.Put(route, object.Put)
	self.Delete(route, object.Delete)
}

// runHooks run each hooks from the Router.
func (self *Router) runHooks(w ResponseWriter, r *Request) *Error {
	for _, hook := range self.hooks {
		if err := hook(w, r); err != nil {
			return err
		}
	}
	return nil
}

// Get registers a pattern with a handler for GET requests.
func (self *Router) Get(route string, h HandlerFunc) *mux.Route {
	return self.Add("GET", route, createHandler(h))
}

// Post registers a pattern with a handler for POST requests.
func (self *Router) Post(route string, h HandlerFunc) *mux.Route {
	return self.Add("POST", route, createHandler(h))
}

// Delete registers a pattern with a handler for DELETE requests.
func (self *Router) Delete(route string, h HandlerFunc) *mux.Route {
	return self.Add("DELETE", route, createHandler(h))
}

// Put registers a pattern with a handler for PUT requests.
func (self *Router) Put(route string, h HandlerFunc) *mux.Route {
	return self.Add("PUT", route, createHandler(h))
}

// ServeHTTP dispatches the handler registered in the matched route.
// It performs any hooks and add the domain registered in the Router to be allowed for cross-domain requests.
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
