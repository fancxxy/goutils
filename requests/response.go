package requests

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

// Response contains *http.Request and *http.Response
type Response struct {
	Req     *http.Request
	Res     *http.Response
	err     error
	reqBody []byte
	resBody []byte
}

// Bytes assign response body to resBody
func (r *Response) Bytes() ([]byte, error) {
	if r.err != nil {
		return nil, r.err
	}

	if r.resBody != nil {
		return r.resBody, nil
	}
	defer r.Res.Body.Close()

	body, err := ioutil.ReadAll(r.Res.Body)
	if err != nil {
		r.err = err
		return nil, err
	}
	r.resBody = body
	return r.resBody, nil
}

// ToString return response body as string
func (r *Response) ToString() string {
	data, _ := r.Bytes()
	return string(data)
}

// ToJSON unmarshal response body to v
func (r *Response) ToJSON(v interface{}) error {
	data, err := r.Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}

// ToFile download response body to file
func (r *Response) ToFile(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := r.Bytes()
	if err != nil {
		return err
	}

	_, err = file.Write(data)
	return err
}

// ToReader return response body as io.Reader
func (r *Response) ToReader() io.Reader {
	data, _ := r.Bytes()
	return bytes.NewReader(data)
}

// String format http request and response
func (r *Response) String() string {
	writer := bytes.NewBuffer(nil)
	r.Req.Write(writer)
	writer.WriteString("\n==============================\n\n")
	r.Res.Write(writer)
	return writer.String()
}
