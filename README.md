# gintest — Django-style integration testing for Gin + GORM + fx

`pkg/gintest` is a small testing harness that turns multi-page boilerplate
into ten-line integration tests, the way Django's `TestCase` + DRF's
`APIClient` do for Python apps. This repo is the example app it was built
against — `pkg/gintest` is the library; everything under `domain/`, `cmd/`,
and `bootstrap/` is the demo target.

## What you get

- **Per-test DB isolation** — every test runs inside a transaction that
  rolls back on cleanup. Tests never leak state to each other.
- **A fluent HTTP client** — `s.Client.POST(path).JSON(body).Send().Status(201)`
  with JSON path assertions backed by `tidwall/gjson`.
- **Factories** — Django/factory-boy style builders for your GORM models.
- **DB assertions** — `s.AssertCount`, `s.AssertExists`, `s.AssertNotExists`.
- **Outbound HTTP mocking** — `s.HTTP.OnGet(url).Reply(200).JSON(...)`.
- **Auth helpers** — `s.Client.AsUser(u)`, `WithBearer`, `WithBasicAuth`.

The library is fx-aware: it accepts your existing fx modules, decorates the
DB provider to inject the transactional session, and extracts the gin engine
from wherever you keep it.

## Quickstart

The full setup for this app lives in [tests/setup.go](tests/setup.go):

```go
func NewSuite(t *testing.T, opts ...gintest.Option) *gintest.Suite {
    return gintest.New(t,
        gintest.WithModules(bootstrap.CommonModules),
        gintest.WithMigrations(func(db *gorm.DB) error {
            return db.AutoMigrate(&models.ProjectCategory{}, &models.Project{})
        }),
        gintest.WithDBDecorator(func(tx *gorm.DB) *infrastructure.Database {
            return &infrastructure.Database{DB: tx}
        }),
        gintest.WithEngineFrom(func(r *infrastructure.Router) *gin.Engine {
            return r.Engine
        }),
    )
}
```

A full test ([tests/project_endpoint_test.go](tests/project_endpoint_test.go)):

```go
func TestCreateProject(t *testing.T) {
    t.Parallel()
    s := NewSuite(t)

    category := factories.ProjectCategory.Create(t, s.DB)

    s.Client.POST("/api/v1/projects").
        JSON(gin.H{
            "name":                "Gaurav Paudel",
            "endpoint":            "https://github.com/paudelgaurav",
            "project_category_id": category.ID,
        }).
        Send().
        Status(http.StatusCreated).
        JSONPath("data.Name").Equals("Gaurav Paudel").
        JSONPath("data.ID").NotEmpty()

    s.AssertCount(&models.Project{}, 1)
}
```

## Defining factories

Define once, reuse everywhere. The builder receives a per-factory sequence
number for uniqueness.

```go
// tests/factories/factories.go
var ProjectCategory = gintest.NewFactory(func(seq int) models.ProjectCategory {
    return models.ProjectCategory{Name: fmt.Sprintf("Category %d", seq)}
})

// Call sites:
cat := factories.ProjectCategory.Create(t, s.DB)
cats := factories.ProjectCategory.CreateN(t, s.DB, 3)
custom := factories.ProjectCategory.Create(t, s.DB, func(c *models.ProjectCategory) {
    c.Name = "Override"
})
```

## Mocking outbound HTTP

Service code that calls `http.Get` can be intercepted without changing the
implementation:

```go
s.HTTP.OnGet("=~^https://example.com/").Reply(200).JSON(gin.H{"ok": true})

s.Client.GET("/api/v1/projects/ping").Send().Status(200)
```

The mock activates lazily on first stub registration and resets on cleanup.
Because `jarcoal/httpmock` patches the process-global default transport, tests
that register stubs are serialized internally — mark `t.Parallel()` on
HTTP-mock tests only when you've verified no other test races on the
same outbound URL.

## Auth

```go
// At suite construction:
gintest.WithAuthProvider(func(u any) (string, string) {
    user := u.(*models.User)
    return "Authorization", "Bearer " + issueToken(user)
})

// In tests:
s.Client.AsUser(user).GET("/api/v1/me").Send().Status(200)
```

Or, if you just need a token directly:

```go
s.Client.WithBearer(token).GET("/api/v1/me").Send()
s.Client.WithBasicAuth("alice", "hunter2").GET("/admin").Send()
```

## Options reference

- `WithModules(opts ...fx.Option)` — your app's fx modules.
- `WithFxOptions(opts ...fx.Option)` — extra options/overrides for the test graph.
- `WithDBOpener(func() (*gorm.DB, error))` — override the base DB (default: isolated in-memory SQLite).
- `WithMigrations(func(*gorm.DB) error)` — run before the per-test transaction begins.
- `WithDBDecorator[T any](func(*gorm.DB) T)` — wrap the tx DB into your app's DB type.
- `WithEngineFrom[T any](func(T) *gin.Engine)` — tell the library where the engine lives.
- `WithAuthProvider(func(any) (name, value string))` — drives `Client.AsUser`.
- `WithSilentFxLogs()` — suppress fx lifecycle logging.

## Files at a glance

- [pkg/gintest/suite.go](pkg/gintest/suite.go) — `Suite`, `New`, lifecycle.
- [pkg/gintest/options.go](pkg/gintest/options.go) — public option constructors.
- [pkg/gintest/client.go](pkg/gintest/client.go) — fluent HTTP client.
- [pkg/gintest/response.go](pkg/gintest/response.go) — `Response`, `Value`, JSONPath assertions.
- [pkg/gintest/factory.go](pkg/gintest/factory.go) — generic `Factory[T]`.
- [pkg/gintest/assert.go](pkg/gintest/assert.go) — DB count / existence helpers.
- [pkg/gintest/httpmock.go](pkg/gintest/httpmock.go) — outbound HTTP mocking.
- [pkg/gintest/auth.go](pkg/gintest/auth.go) — `AsUser`, `WithBearer`, `WithBasicAuth`.
- [pkg/gintest/txdb.go](pkg/gintest/txdb.go) — default in-memory SQLite opener.
