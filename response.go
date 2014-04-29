package http

import (
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

func (self ResponseWriter) AddCORSHeaders(domain string) {
	self.Header().Add("Access-Control-Allow-Origin", domain)
	self.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	self.Header().Add("Access-Control-Allow-Headers", "content-Type, authorization, accept")
	self.Header().Add("Access-Control-Max-Age", "604800")
	self.Header().Add("Access-Control-Allow-Credentials", "true")
}

//  Marshal a single key / value JSON and write it
func (self ResponseWriter) WriteSingleStringJSON(key, value string) {
	if debugMode {
		self.Write([]byte(fmt.Sprintf("{\"%v\":\"%v\"}", key, value)))
	} else {
		self.Write([]byte(fmt.Sprintf("{\n  \"%v\": \"%v\"\n}", key, value)))
	}
}

// Marshal the Object and write it
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

// Send an error using http.Error
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

func (self *Response) debug(str string, values ...interface{}) {
	fmt.Printf("[%v]: %v\n", self.Response.Request.Referer(), fmt.Sprintf("%v%v", str, values))
}

// Returns an the JSON object from the body
func (self *Response) GetJSONObject(object interface{}) (interface{}, error) {
	self.Body.Close()
	body, err := ioutil.ReadAll(self.Body)
	if err != nil {
		return nil, err
	}
	if debugMode {
		self.debug("Body: \"%v\"", string(body))
	}
	return object, json.Unmarshal(body, &object)
}
