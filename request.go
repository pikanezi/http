package http

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

type Request struct {
	*http.Request
}

func createRequest(r *http.Request) *Request {
	return &Request{r}
}

func (self *Request) debug(str string, values ...interface{}) {
	fmt.Printf("[%v]: %v\n", self.RequestURI, fmt.Sprintf(str, values...))
}

// Upload the file, create a new file to the given path (for example "/tmp/")
func (self *Request) UploadAndGetFile(key, pathFile string) (*os.File, error) {
	if debugMode {
		self.debug("Trying to get file from the key \"%v\"", key)
	}
	fileMultiPart, fileHeader, err := self.FormFile(key)
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

// Same as UploadAndGetFile but returns the absolute path of the file
func (self *Request) UploadAndGetAbsolutePath(key, pathFile string) (string, error) {
	file, err := self.UploadAndGetFile(key, pathFile)
	if err != nil {
		return "", err
	}
	defer file.Close()
	absName, err := filepath.Abs(file.Name())
	if err != nil {
		return "", err
	}
	if debugMode {
		self.debug("Got file \"%v\"", absName)
	}
	return absName, nil
}

func (self *Request) getBody() ([]byte, error) {
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

// Returns the JSON object from the body
func (self *Request) GetAndReturnJSONObject(object interface{}) (interface{}, error) {
	body, err := self.getBody()
	if err != nil {
		return nil, err
	}
	return object, json.Unmarshal(body, &object)
}

// Just call json.Unmarshal to the body and put it in the object
func (self *Request) GetJSONObject(object interface{}) error {
	body, err := self.getBody()
	if err != nil {
		return err
	}
	return json.Unmarshal(body, &object)
}
