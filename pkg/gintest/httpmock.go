package gintest

import (
	"sync"
	"testing"

	"github.com/jarcoal/httpmock"
)

// HTTPMock wraps jarcoal/httpmock. It activates lazily — only when a stub is
// first registered — so suites that don't mock outbound HTTP impose no global
// state. Activation is serialized via a package-level mutex.
//
// CAVEAT: jarcoal/httpmock patches http.DefaultTransport, which is process
// global. Tests that register stubs are not safe to run in parallel with
// each other; mark them t.Parallel() at your own risk. Tests that do not
// touch HTTPMock are unaffected.
type HTTPMock struct {
	t      *testing.T
	active bool
}

var httpmockMu sync.Mutex

func newHTTPMock(t *testing.T) *HTTPMock {
	t.Helper()
	return &HTTPMock{t: t}
}

func (m *HTTPMock) ensureActive() {
	if m.active {
		return
	}
	httpmockMu.Lock()
	httpmock.Activate()
	m.active = true
}

func (m *HTTPMock) deactivate() {
	if m == nil || !m.active {
		return
	}
	httpmock.DeactivateAndReset()
	m.active = false
	httpmockMu.Unlock()
}

// OnGet registers a stub for GET requests to url. Returns a stub builder so
// you can describe the response.
func (m *HTTPMock) OnGet(url string) *Stub { return m.on("GET", url) }

// OnPost registers a stub for POST requests to url.
func (m *HTTPMock) OnPost(url string) *Stub { return m.on("POST", url) }

// OnPut registers a stub for PUT requests to url.
func (m *HTTPMock) OnPut(url string) *Stub { return m.on("PUT", url) }

// OnPatch registers a stub for PATCH requests to url.
func (m *HTTPMock) OnPatch(url string) *Stub { return m.on("PATCH", url) }

// OnDelete registers a stub for DELETE requests to url.
func (m *HTTPMock) OnDelete(url string) *Stub { return m.on("DELETE", url) }

// On registers a stub for the given method and URL. Use this for less-common
// verbs not covered by the shortcuts.
func (m *HTTPMock) On(method, url string) *Stub { return m.on(method, url) }

// Reset clears all currently registered stubs. The mock remains active.
func (m *HTTPMock) Reset() { httpmock.Reset() }

// CallCount returns the total number of stubbed calls observed (any method).
func (m *HTTPMock) CallCount() int { return httpmock.GetTotalCallCount() }

func (m *HTTPMock) on(method, url string) *Stub {
	m.ensureActive()
	return &Stub{t: m.t, method: method, url: url, status: 200}
}

// Stub describes a stubbed outbound HTTP response.
type Stub struct {
	t       *testing.T
	method  string
	url     string
	status  int
	body    string
	headers map[string]string
}

// Reply sets the response status code.
func (s *Stub) Reply(status int) *Stub {
	s.status = status
	return s
}

// BodyString sets the response body to a literal string.
func (s *Stub) BodyString(body string) *Stub {
	s.body = body
	s.register(httpmock.NewStringResponder(s.status, body))
	return s
}

// JSON sets the response body to the JSON encoding of v.
func (s *Stub) JSON(v any) *Stub {
	resp, err := httpmock.NewJsonResponder(s.status, v)
	if err != nil {
		s.t.Fatalf("gintest: build JSON responder: %v", err)
	}
	s.register(resp)
	return s
}

// Empty registers an empty-body responder at the configured status.
func (s *Stub) Empty() *Stub {
	s.register(httpmock.NewStringResponder(s.status, ""))
	return s
}

func (s *Stub) register(resp httpmock.Responder) {
	httpmock.RegisterResponder(s.method, s.url, resp)
}
