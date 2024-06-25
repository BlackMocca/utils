package resty

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
)

const (
	methodGet    = "GET"
	methodPost   = "POST"
	methodPUT    = "PUT"
	methodDELETE = "DELETE"
)

const (
	respTypeString = "string"
	respTypeJSON   = "json"
	respTypeRaw    = "raw"
)

type Client struct {
	httptest   *httptest.Server
	method     string
	path       string
	statusCode int
	header     http.Header
	respHeader http.Header
	respBody   interface{}
	respType   string
	handler    http.HandlerFunc
}

func New() *Client {
	return &Client{header: make(http.Header), respHeader: make(http.Header)}
}

func (c *Client) getHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if c.method != r.Method {
			panic("method request does not match with register method")
		}

		if c.path != r.URL.Path {
			panic("path request does not match with register path")
		}
		/* set status code */
		w.WriteHeader(c.statusCode)
		/* set header */
		if len(c.header) > 0 {
			for key, val := range c.header {
				if len(val) > 0 {
					for _, v := range val {
						w.Header().Add(key, v)
						c.respHeader.Add(key, v)
					}
				}
			}
		}

		/* set body */
		switch c.respType {
		case respTypeString:
			var buf bytes.Buffer
			if _, err := io.Copy(&buf, bytes.NewReader([]byte(c.respBody.(string)))); err != nil {
				panic(fmt.Sprintf("fail to copy body"))
			}
			w.Write(buf.Bytes())
			break
		case respTypeJSON:
			var buf = bytes.NewBuffer(c.respBody.([]byte))
			w.Write(buf.Bytes())
			break
		case respTypeRaw:
			var buf bytes.Buffer
			if _, err := io.Copy(&buf, c.respBody.(io.ReadCloser)); err != nil {
				panic(fmt.Sprintf("fail to copy body"))
			}
			w.Write(buf.Bytes())
			break
		}
	})
}

func (c *Client) Action() *httptest.Server {
	ts := httptest.NewUnstartedServer(c.handler)
	ts.Start()

	return ts
}

/* set status */
func (c *Client) Reply(status int) *Client {
	c.statusCode = status
	return c
}

func (c *Client) AddHeader(key, val string) *Client {
	c.header.Add(key, val)
	return c
}

/* request body */
func (c *Client) RawBody(val io.ReadCloser) *Client {
	c.respType = respTypeRaw
	c.respBody = val
	c.handler = c.getHandler()
	return c
}

func (c *Client) BodyString(val string) *Client {
	c.respType = respTypeString
	c.respBody = val
	c.handler = c.getHandler()
	return c
}

func (c *Client) BodyJSON(val interface{}) *Client {
	bu, err := json.Marshal(val)
	if err != nil {
		panic(fmt.Sprintf("fail to parse json: %s", err.Error()))
	}
	c.respType = respTypeJSON
	c.respBody = bu
	c.handler = c.getHandler()
	return c
}

/* request action */
func (c *Client) Get(path string) *Client {
	c.method = methodGet
	path = strings.Trim(path, "/")
	c.path = fmt.Sprintf("/%s", path)
	return c
}

func (c *Client) Post(path string) *Client {
	c.method = methodPost
	path = strings.Trim(path, "/")
	c.path = fmt.Sprintf("/%s", path)
	return c
}

func (c *Client) Put(path string) *Client {
	c.method = methodPUT
	path = strings.Trim(path, "/")
	c.path = fmt.Sprintf("/%s", path)
	return c
}

func (c *Client) Delete(path string) *Client {
	c.method = methodDELETE
	path = strings.Trim(path, "/")
	c.path = fmt.Sprintf("/%s", path)
	return c
}

/* get Attribute */
func (c *Client) GetServer() *httptest.Server {
	return c.httptest
}
func (c *Client) GetPath() string {
	return c.path
}
func (c *Client) GetMethod() string {
	return c.method
}
func (c *Client) GetStatusCode() int {
	return c.statusCode
}
func (c *Client) GetHeader() http.Header {
	return c.respHeader
}
func (c *Client) GetBody() interface{} {
	return c.respBody
}

func (c *Client) SetHandlerFunc(function http.HandlerFunc) {
	c.handler = function
}

func (c *Client) StartTLS() {
	c.httptest.StartTLS()
}
func (c *Client) Start() {
	c.httptest.Start()
}
func (c *Client) Close() {
	c.httptest.Close()
}
func (c *Client) CloseClientConnections() {
	c.httptest.CloseClientConnections()
}
