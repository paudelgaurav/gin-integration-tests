package gintest

import (
	"encoding/json"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/tidwall/gjson"
)

// Response wraps an httptest.ResponseRecorder with assertion helpers.
// Every method calls t.Helper() and fails the test via t.Fatalf on mismatch.
type Response struct {
	t   *testing.T
	rec *httptest.ResponseRecorder
}

// Code returns the HTTP status code.
func (r *Response) Code() int { return r.rec.Code }

// Body returns the response body as bytes.
func (r *Response) Body() []byte { return r.rec.Body.Bytes() }

// BodyString returns the response body as a string.
func (r *Response) BodyString() string { return r.rec.Body.String() }

// HeaderValue returns the value of a response header.
func (r *Response) HeaderValue(key string) string { return r.rec.Header().Get(key) }

// Status asserts the response has the expected status code.
func (r *Response) Status(expected int) *Response {
	r.t.Helper()
	if r.rec.Code != expected {
		r.t.Fatalf("expected status %d, got %d\nbody: %s", expected, r.rec.Code, r.rec.Body.String())
	}
	return r
}

// HeaderEquals asserts a response header equals the expected value.
func (r *Response) HeaderEquals(key, expected string) *Response {
	r.t.Helper()
	if got := r.rec.Header().Get(key); got != expected {
		r.t.Fatalf("expected header %q=%q, got %q", key, expected, got)
	}
	return r
}

// BodyContains asserts the response body contains the substring s.
func (r *Response) BodyContains(s string) *Response {
	r.t.Helper()
	if !strings.Contains(r.rec.Body.String(), s) {
		r.t.Fatalf("expected body to contain %q\nbody: %s", s, r.rec.Body.String())
	}
	return r
}

// JSONPath returns a Value at the given gjson path (e.g. "data.items.0.name")
// for further assertions.
func (r *Response) JSONPath(path string) *Value {
	r.t.Helper()
	return &Value{
		t:    r.t,
		path: path,
		res:  gjson.GetBytes(r.rec.Body.Bytes(), path),
		resp: r,
	}
}

// DecodeJSON unmarshals the response body into the target.
func (r *Response) DecodeJSON(target any) *Response {
	r.t.Helper()
	if err := json.Unmarshal(r.rec.Body.Bytes(), target); err != nil {
		r.t.Fatalf("gintest: decode JSON: %v\nbody: %s", err, r.rec.Body.String())
	}
	return r
}

// Value is a single JSON-path lookup result with chainable assertions.
type Value struct {
	t    *testing.T
	path string
	res  gjson.Result
	resp *Response
}

// End returns the underlying Response so further response-level assertions
// or JSON path queries can be chained. Example:
//
//	res.JSONPath("data.id").NotEmpty().End().JSONPath("data.name").Equals("Foo")
func (v *Value) End() *Response { return v.resp }

// JSONPath is a shortcut for v.End().JSONPath(path), enabling fluent chaining
// across multiple paths on the same response.
func (v *Value) JSONPath(path string) *Value { return v.resp.JSONPath(path) }

// Exists asserts the path resolved to a value (including JSON null).
func (v *Value) Exists() *Value {
	v.t.Helper()
	if !v.res.Exists() {
		v.t.Fatalf("expected JSON path %q to exist", v.path)
	}
	return v
}

// NotEmpty asserts the path resolved to a non-empty, non-null value.
func (v *Value) NotEmpty() *Value {
	v.t.Helper()
	if !v.res.Exists() || v.res.Type == gjson.Null {
		v.t.Fatalf("expected JSON path %q to be non-empty, got %v", v.path, v.res.Raw)
	}
	switch v.res.Type {
	case gjson.String:
		if v.res.Str == "" {
			v.t.Fatalf("expected JSON path %q to be non-empty string", v.path)
		}
	case gjson.Number:
		// zero is a legitimate non-empty number; nothing to check
	case gjson.JSON:
		if v.res.Raw == "[]" || v.res.Raw == "{}" {
			v.t.Fatalf("expected JSON path %q to be non-empty, got %s", v.path, v.res.Raw)
		}
	}
	return v
}

// Equals asserts the JSON value at the path equals expected. Numbers, strings,
// bools, and JSON-comparable types are supported. For complex types, expected
// is JSON-marshaled and compared structurally.
func (v *Value) Equals(expected any) *Value {
	v.t.Helper()
	if !v.res.Exists() {
		v.t.Fatalf("expected JSON path %q to equal %v, but path missing", v.path, expected)
	}

	switch want := expected.(type) {
	case string:
		if v.res.String() != want {
			v.t.Fatalf("JSON path %q: expected %q, got %q", v.path, want, v.res.String())
		}
	case bool:
		if v.res.Bool() != want {
			v.t.Fatalf("JSON path %q: expected %v, got %v", v.path, want, v.res.Bool())
		}
	case int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64,
		float32, float64:
		got := v.res.Float()
		wantF := reflect.ValueOf(expected).Convert(reflect.TypeOf(float64(0))).Float()
		if got != wantF {
			v.t.Fatalf("JSON path %q: expected %v, got %v", v.path, want, got)
		}
	default:
		// Structural compare via JSON round-trip.
		wantJSON, err := json.Marshal(expected)
		if err != nil {
			v.t.Fatalf("gintest: marshal expected: %v", err)
		}
		if !jsonEqual([]byte(v.res.Raw), wantJSON) {
			v.t.Fatalf("JSON path %q: expected %s, got %s", v.path, wantJSON, v.res.Raw)
		}
	}
	return v
}

// Len asserts the array or object at the path has n elements.
func (v *Value) Len(n int) *Value {
	v.t.Helper()
	if !v.res.Exists() {
		v.t.Fatalf("expected JSON path %q to have length %d, but path missing", v.path, n)
	}
	var got int
	if v.res.IsArray() {
		got = len(v.res.Array())
	} else if v.res.IsObject() {
		got = len(v.res.Map())
	} else {
		v.t.Fatalf("JSON path %q: value is not array/object: %s", v.path, v.res.Raw)
	}
	if got != n {
		v.t.Fatalf("JSON path %q: expected length %d, got %d", v.path, n, got)
	}
	return v
}

// String returns the underlying string value (or the raw JSON for non-string types).
func (v *Value) String() string { return v.res.String() }

// Int returns the underlying integer value.
func (v *Value) Int() int64 { return v.res.Int() }

// Raw returns the raw JSON at the path.
func (v *Value) Raw() string { return v.res.Raw }

func jsonEqual(a, b []byte) bool {
	var ax, bx any
	if err := json.Unmarshal(a, &ax); err != nil {
		return false
	}
	if err := json.Unmarshal(b, &bx); err != nil {
		return false
	}
	return reflect.DeepEqual(ax, bx)
}
