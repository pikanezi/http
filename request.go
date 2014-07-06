package http

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// Request represents an HTTP request received by a server or to be sent by a client.
type Request struct {
	*http.Request
}

// Create a new Request from a classic Request.
func CreateRequest(r *http.Request) *Request {
	return &Request{r}
}

// NewRequest returns a new Request given a method, URL, and optional body.
func NewRequest(method, url string, body io.Reader) (*Request, error) {
	r, err := http.NewRequest(method, url, body)
	return createRequest(r), err
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

// UploadAndGetFile upload the file, create a new file to the given path (for example "/tmp/").
func (req *Request) UploadAndGetFile(key, pathFile string) (*os.File, error) {
	if debugMode {
		req.debug("Trying to get file from the key \"%v\"", key)
	}
	fileMultiPart, fileHeader, err := req.FormFile(key)
	if err != nil {
		return nil, err
	}
	file, err := os.Create(fmt.Sprintf("%v%v", pathFile, fileHeader.Filename))
	if err != nil {
		return nil, err
	}
	if _, err := io.Copy(file, fileMultiPart); err != nil {
		return nil, err
	}
	return file, nil
}

// UploadAndGetAbsolutePath is the same as UploadAndGetFile but returns the absolute path of the file.
func (req *Request) UploadAndGetAbsolutePath(key, pathFile string) (string, error) {
	file, err := req.UploadAndGetFile(key, pathFile)
	if err != nil {
		return "", err
	}
	defer file.Close()
	absName, err := filepath.Abs(file.Name())
	if err != nil {
		return "", err
	}
	if debugMode {
		req.debug("Got file \"%v\"", absName)
	}
	return absName, nil
}

// UploadFileName uploads the file from the request and save it in the given fileName.
func (req *Request) UploadFileName(key, fileName string) (*os.File, error) {
	if debugMode {
		req.debug("Trying to get file from the key \"%v\"", key)
	}
	fileMultiPart, _, err := req.FormFile(key)
	if err != nil {
		return nil, err
	}
	file, err := os.Create(fileName)
	if err != nil {
		return nil, err
	}
	if _, err := io.Copy(file, fileMultiPart); err != nil {
		return nil, err
	}
	return file, nil
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

// URLParam returns an URL param.
// It is the same as calling request.Url.Query().Get(":key").
func (req *Request) URLParam(key string) string {
	return req.Request.URL.Query().Get(":" + key)
}
