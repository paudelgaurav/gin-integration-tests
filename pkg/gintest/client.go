package gintest

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
)

// Client is a thin wrapper that issues requests against the suite's gin.Engine.
// It carries sticky headers (set via WithHeader / AsUser) so authenticated
// flows don't repeat themselves in every test.
type Client struct {
	suite   *Suite
	headers http.Header
}

func newClient(s *Suite) *Client {
	return &Client{suite: s, headers: http.Header{}}
}

// Clone returns a copy of the client with independent sticky headers.
// Useful when one test wants to make calls "as different users".
func (c *Client) Clone() *Client {
	cp := &Client{suite: c.suite, headers: http.Header{}}
	for k, vs := range c.headers {
		cp.headers[k] = append([]string(nil), vs...)
	}
	return cp
}

// WithHeader sets a sticky header applied to every subsequent request.
// Returns the client for chaining.
func (c *Client) WithHeader(key, value string) *Client {
	c.headers.Set(key, value)
	return c
}

// GET starts building a GET request.
func (c *Client) GET(path string) *Request { return c.req(http.MethodGet, path) }

// POST starts building a POST request.
func (c *Client) POST(path string) *Request { return c.req(http.MethodPost, path) }

// PUT starts building a PUT request.
func (c *Client) PUT(path string) *Request { return c.req(http.MethodPut, path) }

// PATCH starts building a PATCH request.
func (c *Client) PATCH(path string) *Request { return c.req(http.MethodPatch, path) }

// DELETE starts building a DELETE request.
func (c *Client) DELETE(path string) *Request { return c.req(http.MethodDelete, path) }

func (c *Client) req(method, path string) *Request {
	r := &Request{
		client:  c,
		method:  method,
		path:    path,
		headers: http.Header{},
		query:   url.Values{},
	}
	for k, vs := range c.headers {
		r.headers[k] = append([]string(nil), vs...)
	}
	return r
}

// Request is a builder for a single HTTP call.
type Request struct {
	client  *Client
	method  string
	path    string
	body    io.Reader
	headers http.Header
	query   url.Values
}

// JSON sets the request body to the JSON encoding of v and sets the
// Content-Type header to application/json.
func (r *Request) JSON(v any) *Request {
	buf, err := json.Marshal(v)
	if err != nil {
		r.client.suite.T.Fatalf("gintest: marshal JSON body: %v", err)
	}
	r.body = bytes.NewReader(buf)
	if r.headers.Get("Content-Type") == "" {
		r.headers.Set("Content-Type", "application/json")
	}
	return r
}

// Body sets a raw request body. Caller is responsible for Content-Type.
func (r *Request) Body(reader io.Reader) *Request {
	r.body = reader
	return r
}

// Form sets a urlencoded form body.
func (r *Request) Form(values url.Values) *Request {
	r.body = strings.NewReader(values.Encode())
	if r.headers.Get("Content-Type") == "" {
		r.headers.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	return r
}

// Header overrides a single header on this request.
func (r *Request) Header(key, value string) *Request {
	r.headers.Set(key, value)
	return r
}

// Query adds a query-string parameter.
func (r *Request) Query(key, value string) *Request {
	r.query.Add(key, value)
	return r
}

// Send executes the request through the suite's gin.Engine and returns a
// Response for assertions.
func (r *Request) Send() *Response {
	r.client.suite.T.Helper()

	target := r.path
	if len(r.query) > 0 {
		sep := "?"
		if strings.Contains(target, "?") {
			sep = "&"
		}
		target = target + sep + r.query.Encode()
	}

	req := httptest.NewRequest(r.method, target, r.body)
	for k, vs := range r.headers {
		req.Header[k] = vs
	}

	rec := httptest.NewRecorder()
	r.client.suite.Engine.ServeHTTP(rec, req)

	return &Response{
		t:   r.client.suite.T,
		rec: rec,
	}
}
