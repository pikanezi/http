Http
====

Override of net/http and gorilla/pat golang library to use JSON with Request and ResponseWritter in a more easy way.

Its purpose is to create web server sending and receiving JSON data in few lines.

See documentation of gorilla/pat here : http://www.gorillatoolkit.org/pkg/pat

---

Full Example
====

```go

import (
	"github.com/pikanezi/http"
	"io"
	"log"
)

const (
	errorWrongObject = iota << 2
)

type Object struct {
	SomeField string `json:"someField,omitempty"`
}

func HelloWorldHandler(w http.ResponseWriter, r *http.Request) *http.Error {
	object := &Object{"Hello World"}
	w.WriteJSON(object)
	return nil
}

func AdminHandler(w http.ResponseWriter, r *http.Request) *http.Error {
	w.WriteJSON("Hello Admin!")
	return nil
}

func CheckAdminInterceptor(w http.ResponseWriter, r *http.Request) *http.Error {
	// Check if user is an Admin
	return nil
}


func AddObjectHandler(w http.ResponseWriter, r *http.Request) *http.Error {
	object := &Object{}
	if err := r.GetJSONObject(&object); err != nil {
		return http.NewError(err, 500)
	}
	// handle object
	return nil
}

// GetFileReader works only for single file reader.
func GetFileHandler(w http.ResponseWriter, r *http.Request) *http.Error {
	reader, err := r.GetFileReader("file")
	// Handle file Reader
	return nil
}

func MultiFileHandler(w http.ResponseWriter, r *http.Request) *http.Error {
	if err := r.ForEachFile("files", func (index int, reader io.Reader) {
		// Do something with each file Reader
	}); err != nil {
		return http.NewError(err, 500)
	}
	return nil
}

func main() {
	r := http.NewRouter()
    
	// Set the custom headers to be set before the Handler handle the request.
	// It is useful for handling the CORS.
	r.SetCustomHeader(http.Header{
		"Access-Control-Allow-Origin":      domainCORS,
		"Access-Control-Allow-Methods":     "POST, GET, OPTIONS, PUT, DELETE",
		"Access-Control-Allow-Headers":     "Content-Type, Authorization, Accept, x-api-key",
		"Access-Control-Allow-Max-Age":     "604800",
		"Access-Control-Allow-Credentials": "true",
	})
	
	r.Get("/hello/world", HelloWorldHandler).Before(CheckAdminInterceptor)
	r.Post("/object", AddObjectHandler)
    r.Post("/file", GetFileHandler)
    
	log.Fatal(http.ListenAndServe(":8080", r)
}

```

License
====

The MIT License (MIT)

Copyright (c) 2014 Vincent Nëel

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
