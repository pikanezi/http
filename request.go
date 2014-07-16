package http

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"encoding/xml"
)

// Request represents an HTTP request received by a server or to be sent by a client.
type Request struct {
	*http.Request
}

// CreateRequest a new Request from a classic Request.
func CreateRequest(r *http.Request) *Request {
	return &Request{r}
}

// NewRequest returns a new Request given a method, URL, and optional body.
func NewRequest(method, url string, body io.Reader) (*Request, error) {
	r, err := http.NewRequest(method, url, body)
	return CreateRequest(r), err
}

func (req *Request) debug(str string, values ...interface{}) {
	log.Printf("[%v]: %v\n", req.Request.RequestURI, fmt.Sprintf(str, values...))
}

// GetFileReader get the multiform body and returns it as a Reader.
func (req *Request) GetFileReader(key string) (io.Reader, error) {
	if debugMode {
		req.debug("Trying to get file from the key \"%v\"", key)
	}
	fileMultiPart, _, err := req.FormFile(key)
	if err != nil {
		return nil, err
	}
	return fileMultiPart, nil
}

// ForEachFileReader calls the function f with every io.Reader contained in the provided form key.
// After every call to f, the file is close.
func (req *Request) ForEachFileReader(key string, f func(int, io.Reader) error) error {
	if err := req.ParseMultipartForm(32 << 2); err != nil {
		return err
	}
	if req.MultipartForm != nil && req.MultipartForm.File[key] != nil {
		fileHeaders := req.MultipartForm.File[key]
		for index, fileheader := range fileHeaders {
			file, err := fileheader.Open()
			if err != nil {
				return err
			}
			if err := f(index, file); err != nil {
				return err
			}
			if err := file.Close(); err != nil {
				return err
			}
		}
	}
	return nil
}

// ForEachFileHeader calls the function f with every multipart.FileHeader contained in the provided form key.
// Opening and closing the file is the responsibility of the user.
func (req *Request) ForEachFileHeader(key string, f func(int, *multipart.FileHeader) error) error {
	if err := req.ParseMultipartForm(32 << 2); err != nil {
		return err
	}
	if req.MultipartForm != nil && req.MultipartForm.File[key] != nil {
		fileHeaders := req.MultipartForm.File[key]
		for index, fileheader := range fileHeaders {
			if err := f(index, fileheader); err != nil {
				return err
			}
		}
	}
	return nil
}

func (req *Request) getBody() ([]byte, error) {
	defer req.Body.Close()
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}
	if debugMode {
		req.debug("Body: \"%v\"", string(body))
	}
	return body, err
}

// GetAndReturnJSONObject returns the JSON object from the body.
func (req *Request) GetAndReturnJSONObject(object interface{}) (interface{}, error) {
	body, err := req.getBody()
	if err != nil {
		return nil, err
	}
	return object, json.Unmarshal(body, &object)
}

// GetJSONObject just call json.Unmarshal to the body and put it in the object.
func (req *Request) GetJSONObject(object interface{}) error {
	body, err := req.getBody()
	if err != nil {
		return err
	}
	return json.Unmarshal(body, &object)
}

// GetXMLObject jus call xml.Unmarshal to the body and put it in the object.
func (req *Request) GetXMLObject(object interface{}) error {
	body, err := req.getBody()
	if err != nil {
		return err
	}
	return xml.Unmarshal(body, &object)
}

// URLParam returns an URL param.
// It is the same as calling request.Url.Query().Get(":key").
func (req *Request) URLParam(key string) string {
	return req.Request.URL.Query().Get(":" + key)
}
