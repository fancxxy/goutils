package requests

import (
	"compress/gzip"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"
)

var defaultClient *http.Client

func init() {
	jar, _ := cookiejar.New(nil)
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	defaultClient = &http.Client{
		Jar:       jar,
		Transport: transport,
		Timeout:   2 * time.Minute,
	}
}

// Client contains request client
type Client struct {
	client *http.Client
}

// Header represents http request header
type Header map[string]string

// Query represents http request query param
type Query = url.Values

func (r *Client) do(method, rawurl string, params ...interface{}) (*Response, error) {
	var (
		req = &http.Request{
			Method:     method,
			Header:     make(http.Header),
			Proto:      "HTTP/1.1",
			ProtoMajor: 1,
			ProtoMinor: 1,
		}
		query, form Query
		response    = &Response{Req: req}
	)

	for _, param := range params {
		switch m := param.(type) {
		case Header:
			for key, value := range m {
				req.Header.Add(key, value)
			}
		case Query:
			switch method {
			case "GET":
				query = m
			case "POST":
				form = m
			}
		}
	}

	if form != nil {
		body := form.Encode()
		response.reqBody = []byte(body)
		req.Body = ioutil.NopCloser(strings.NewReader(body))
		if req.Header.Get("Content-Type") == "" {
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
		}
	}

	if query != nil {
		rawurl = rawurl + "?" + query.Encode()
	}
	u, err := url.Parse(rawurl)
	if err != nil {
		return nil, err
	}
	req.URL = u

	if host := req.Header.Get("Host"); host != "" {
		req.Host = host
	}

	res, err := r.client.Do(req)
	if err != nil {
		return nil, err
	}

	response.Res = res

	if res.Header.Get("Content-Encoding") == "gzip" && req.Header.Get("Accept-Encoding") != "" {
		body, err := gzip.NewReader(res.Body)
		if err != nil {
			return nil, err
		}
		res.Body = body
	}

	return response, nil
}

// New create *Client
func New(clients ...*http.Client) *Client {
	var client *http.Client
	if len(clients) == 0 {
		client = defaultClient
	} else {
		client = clients[0]
	}

	req := &Client{client: client}
	return req
}

// Get request
func (r *Client) Get(url string, v ...interface{}) (*Response, error) {
	return r.do("GET", url, v...)
}

// Post request
func (r *Client) Post(url string, v ...interface{}) (*Response, error) {
	return r.do("POST", url, v...)
}
