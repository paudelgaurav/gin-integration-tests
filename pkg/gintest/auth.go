package gintest

import "encoding/base64"

// AsUser returns a new Client whose subsequent requests carry the auth
// header produced by the WithAuthProvider callback registered on the suite.
// If no provider is registered, the test fails.
//
// Provider semantics: given the user-like value, return (headerName, headerValue).
// Returning an empty headerName disables injection. This lets callers represent
// "no user / unauthenticated" without branching at the call site.
func (c *Client) AsUser(user any) *Client {
	c.suite.T.Helper()
	if c.suite.authProvider == nil {
		c.suite.T.Fatalf("gintest: AsUser called but no auth provider configured; use gintest.WithAuthProvider")
	}
	cp := c.Clone()
	name, value := c.suite.authProvider(user)
	if name != "" {
		cp.headers.Set(name, value)
	}
	return cp
}

// WithBearer returns a new Client that sends Authorization: Bearer <token>
// on every request. Convenience for token-based auth without needing a
// provider.
func (c *Client) WithBearer(token string) *Client {
	cp := c.Clone()
	cp.headers.Set("Authorization", "Bearer "+token)
	return cp
}

// WithBasicAuth returns a new Client that sends HTTP Basic auth.
func (c *Client) WithBasicAuth(username, password string) *Client {
	cp := c.Clone()
	cp.headers.Set("Authorization", basicAuthHeader(username, password))
	return cp
}

func basicAuthHeader(user, pass string) string {
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(user+":"+pass))
}
