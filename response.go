package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"encoding/xml"
)

// A ResponseWriter interface is used by an HTTP handler to construct an HTTP response.
type ResponseWriter struct {
	http.ResponseWriter
	http.Hijacker
}

// CreateResponseWriter create a new ResponseWriter from a classic ResponseWriter.
func CreateResponseWriter(r http.ResponseWriter) ResponseWriter {
	return ResponseWriter{r, r.(http.Hijacker)}
}

func (rw ResponseWriter) addCustomPreHeader(headers Header) {
	for key, value := range headers {
		rw.Header().Add(key, value)
	}
}

// WriteSingleStringJSON marshal a single key / value JSON and write it.
func (rw ResponseWriter) WriteSingleStringJSON(key, value string) {
	rw.Header().Add("Content-Type", "application/json")
	if debugMode {
		rw.Write([]byte(fmt.Sprintf("{\n  \"%v\": \"%v\"\n}", key, value)))
	} else {
		rw.Write([]byte(fmt.Sprintf("{\"%v\":\"%v\"}", key, value)))
	}
}

// WriteJSON marshal the Object and write it.
func (rw ResponseWriter) WriteJSON(object interface{}) error {
	rw.Header().Add("Content-Type", "application/json")
	if debugMode {
		js, err := json.MarshalIndent(object, "", "  ")
		if err != nil {
			return err
		}
		rw.Write(js)
	} else {
		js, err := json.Marshal(object)
		if err != nil {
			return err
		}
		rw.Write(js)
	}
	return nil
}

// WriteError send an error using http.Error.
func (rw ResponseWriter) WriteError(customErr *Error) error {
	if debugMode {
		b, err := json.MarshalIndent(customErr, "", "  ")
		if err != nil {
			return err
		}
		http.Error(rw, string(b), customErr.HTTPCode)
	} else {
		b, err := json.Marshal(customErr)
		if err != nil {
			return err
		}
		http.Error(rw, string(b), customErr.HTTPCode)
	}
	return nil
}

// Response represents the response from an HTTP request.
type Response struct {
	*http.Response
}

func createResponse(r *http.Response) *Response { return &Response{r} }

// Get issues a GET request to the given URL and returns the Response.
func Get(url string) (*Response, error) {
	r, err := http.Get(url)
	return createResponse(r), err
}

// PostJSON issues a POST request of type "application/json" (the object will be marshaled as JSON) to the given URL.
func PostJSON(url string, object interface{}) (*Response, error) {
	objectJSON, err := json.Marshal(object)
	if err != nil {
		return nil, err
	}
	r, err := http.Post(url, "application/json", bytes.NewReader(objectJSON))
	return createResponse(r), err
}

// PostForm issues a POST to the specified URL, with data's keys and values URL-encoded as the request body.
func PostForm(url string, data url.Values) (*Response, error) {
	r, err := http.PostForm(url, data)
	return createResponse(r), err
}

func (resp *Response) debug(str string, values ...interface{}) {
	fmt.Printf("[%v]: %v\n", resp.Response.Request.RequestURI, fmt.Sprintf("%v%v", str, values))
}

func (resp *Response) getBody() ([]byte, error) {
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if debugMode {
		resp.debug("Body: %v", string(body))
	}
	return body, err
}

// GetAndReturnJSONObject returns an the JSON object from the body.
func (resp *Response) GetAndReturnJSONObject(object interface{}) (interface{}, error) {
	body, err := resp.getBody()
	if err != nil {
		return nil, err
	}
	return object, json.Unmarshal(body, &object)
}

// GetJSONObject just call json.Unmarshal to the body and put it in the object.
func (resp *Response) GetJSONObject(object interface{}) error {
	body, err := resp.getBody()
	if err != nil {
		return err
	}
	return json.Unmarshal(body, &object)
}

// GetXMLObject jus call xml.Unmarshal to the body and put it in the object.
func (resp *Response) GetXMLObject(object interface{}) error {
	body, err := resp.getBody()
	if err != nil {
		return err
	}
	return xml.Unmarshal(body, &object)
}
