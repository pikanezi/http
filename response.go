package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type ResponseWriter struct {
	http.ResponseWriter
	http.Hijacker
}

func createResponseWriter(r http.ResponseWriter) ResponseWriter {
	return ResponseWriter{r, r.(http.Hijacker)}
}

func (self ResponseWriter) addCORSHeaders(domain string) {
	self.Header().Add("Access-Control-Allow-Origin", domain)
	self.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	self.Header().Add("Access-Control-Allow-Headers", "content-Type, authorization, accept")
	self.Header().Add("Access-Control-Max-Age", "604800")
	self.Header().Add("Access-Control-Allow-Credentials", "true")
}

//  Marshal a single key / value JSON and write it.
func (self ResponseWriter) WriteSingleStringJSON(key, value string) {
	if debugMode {
		self.Write([]byte(fmt.Sprintf("{\n  \"%v\": \"%v\"\n}", key, value)))
	} else {
		self.Write([]byte(fmt.Sprintf("{\"%v\":\"%v\"}", key, value)))
	}
}

// Marshal the Object and write it.
func (self ResponseWriter) WriteJSON(object interface{}) error {
	if debugMode {
		js, err := json.MarshalIndent(object, "", "  ")
		if err != nil {
			return err
		}
		self.Write(js)
	} else {
		js, err := json.Marshal(object)
		if err != nil {
			return err
		}
		self.Write(js)
	}
	return nil
}

// Send an error using http.Error.
func (self ResponseWriter) WriteError(customErr *Error) error {
	if debugMode {
		b, err := json.MarshalIndent(customErr, "", "  ")
		if err != nil {
			return err
		}
		http.Error(self, string(b), customErr.HttpCode)
	} else {
		b, err := json.Marshal(customErr)
		if err != nil {
			return err
		}
		http.Error(self, string(b), customErr.HttpCode)
	}
	return nil
}

type Response struct {
	*http.Response
}

func createResponse(r *http.Response) *Response { return &Response{r} }

// Issues a GET request to the given URL and returns the Response.
func Get(url string) (*Response, error) {
	r, err := http.Get(url)
	return createResponse(r), err
}

// Issues a POST request of type "application/json" (the object will be marshaled as JSON) to the given URL.
func PostJSON(url string, object interface{}) (*Response, error) {
	objectJSON, err := json.Marshal(object)
	if err != nil {
		return nil, err
	}
	r, err := http.Post(url, "application/json", bytes.NewReader(objectJSON))
	return createResponse(r), err
}

// PostForm issues a POST to the specified URL, with data's keys and values URL-encoded as the request body.
func PostForm(url string, data url.values) (*Response, error) {
	r, err := http.PostForm(url, data)
	return createResponse(r), err
}

func (self *Response) debug(str string, values ...interface{}) {
	fmt.Printf("[%v]: %v\n", self.Response.Request.RequestURI, fmt.Sprintf("%v%v", str, values))
}

func (self *Response) getBody() ([]byte, error) {
	defer self.Body.Close()
	body, err := ioutil.ReadAll(self.Body)
	if err != nil {
		return nil, err
	}
	if debugMode {
		self.debug("Body: \"%v\"", string(body))
	}
	return body, err
}

// Returns an the JSON object from the body.
func (self *Response) GetAndReturnJSONObject(object interface{}) (interface{}, error) {
	body, err := self.getBody()
	if err != nil {
		return nil, err
	}
	return object, json.Unmarshal(body, &object)
}

// Just call json.Unmarshal to the body and put it in the object.
func (self *Response) GetJSONObject(object interface{}) error {
	body, err := self.getBody()
	if err != nil {
		return err
	}
	return json.Unmarshal(body, &object)
}
