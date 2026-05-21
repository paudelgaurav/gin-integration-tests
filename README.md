# gintest — Django-style integration testing for Gin + GORM + fx

`pkg/gintest` is a small testing harness that turns multi-page boilerplate
into ten-line integration tests, the way Django's `TestCase` + DRF's
`APIClient` do for Python apps. This repo is the example app it was built
against — `pkg/gintest` is the library; everything under `domain/`, `cmd/`,
and `bootstrap/` is the demo target.

## Install

```bash
go get github.com/paudelgaurav/gin-integration-tests/pkg/gintest
```

Requires Go 1.23+, a Gin app wired with `go.uber.org/fx`, and GORM for
the database layer.

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

## Test database

### Default: isolated in-memory SQLite (zero setup)

If you don't pass `WithDBOpener`, the library opens a fresh in-memory
SQLite database for each call to `gintest.New`, runs your migrations,
and wraps the test in a transaction that rolls back on cleanup.

This is fast and parallel-safe — each suite gets its own DB via a unique
shared-cache DSN. The tradeoff: SQLite SQL is not identical to
Postgres/MySQL (e.g. no `RETURNING`, different JSON operators, weaker
type coercion). For schema-heavy apps you'll want a real DB.

### Postgres / MySQL via testcontainers

Spin up one container for the whole `go test` run; let each test grab a
connection. Sketch:

```go
// tests/dbmain_test.go
var pgDSN string

func TestMain(m *testing.M) {
    ctx := context.Background()
    container, err := postgres.RunContainer(ctx,
        testcontainers.WithImage("postgres:16-alpine"),
        postgres.WithDatabase("test"),
        postgres.WithUsername("test"),
        postgres.WithPassword("test"),
    )
    if err != nil { panic(err) }
    pgDSN, _ = container.ConnectionString(ctx, "sslmode=disable")
    code := m.Run()
    _ = container.Terminate(ctx)
    os.Exit(code)
}

// tests/setup.go
gintest.WithDBOpener(func() (*gorm.DB, error) {
    return gorm.Open(postgres.Open(pgDSN), &gorm.Config{})
}),
gintest.WithMigrations(func(db *gorm.DB) error {
    return db.AutoMigrate(/* your models */)
}),
```

Transaction-per-test isolation still applies, so tests don't see each
other's writes even though they share a database.

> Tip: if `AutoMigrate` is slow against a real DB, run it once in
> `TestMain` against `baseDB` directly, then have `WithDBOpener` return
> the existing connection (no migration in `WithMigrations`). The
> library will still wrap each test in its own tx.

## Mocking outbound HTTP

`s.HTTP` patches the process-global `http.DefaultTransport`, so any code
that calls `http.Get`, `http.DefaultClient`, or constructs a client
without overriding the transport will be intercepted:

```go
s.HTTP.OnGet("=~^https://example.com/").Reply(200).JSON(gin.H{"ok": true})

s.Client.GET("/api/v1/projects/ping").Send().Status(200)
```

The mock activates lazily on first stub registration and resets on cleanup.
Because `jarcoal/httpmock` mutates global state, tests that register stubs
are serialized internally — mark `t.Parallel()` on HTTP-mock tests only
when you've verified no other test races on the same outbound URL.

### `=~` regex prefix

URLs starting with `=~` are treated as regexes — useful for matching any
path under a host:

```go
s.HTTP.OnPost("=~^https://api.stripe.com/v1/").Reply(200).JSON(...)
```

## Mocking external services (S3, Stripe, Redis, …)

There are three shapes you'll run into. The right tool depends on how
the SDK is wired up, not on the vendor.

### A. SDKs that use the default `http.Client` — Stripe, Sendgrid, most REST SDKs

These are intercepted automatically by `s.HTTP` — no code changes needed.

```go
s.HTTP.OnPost("=~^https://api.stripe.com/v1/charges").
    Reply(200).
    JSON(map[string]any{"id": "ch_123", "status": "succeeded"})

s.Client.POST("/api/v1/checkout").JSON(...).Send().Status(200)
```

### B. SDKs with their own HTTP client — AWS SDK v2, GCP, Twilio

These bypass `http.DefaultTransport`, so `s.HTTP` won't see their
traffic. You have two options:

**B1. Inject `httpmock`'s transport into the SDK's client.** AWS SDK v2:

```go
import "github.com/jarcoal/httpmock"

cfg, _ := config.LoadDefaultConfig(ctx,
    config.WithHTTPClient(&http.Client{Transport: httpmock.DefaultTransport}),
)
```

You'd typically swap the AWS config in tests via
`WithFxOptions(fx.Decorate(...))`, then register stubs as usual.

**B2 (preferred). Mock at your own interface.** Most teams wrap S3 behind
an `Uploader` interface and only mock that — see section C below. It
keeps tests aligned with how your code actually uses the SDK rather
than serializing/deserializing fake AWS XML.

### C. Interface-based services — your own Redis client, queue producer, S3 wrapper

This is what fx.Decorate is built for. Provide a fake implementation
and decorate the provider:

```go
// In your app:
type Mailer interface { Send(to, subj, body string) error }
// fx.Provide(func() Mailer { return realMailer{...} })

// In tests:
type fakeMailer struct{ Sent []Email }
func (f *fakeMailer) Send(to, subj, body string) error {
    f.Sent = append(f.Sent, Email{to, subj, body})
    return nil
}

mailer := &fakeMailer{}

s := NewSuite(t,
    gintest.WithFxOptions(fx.Decorate(func(Mailer) Mailer { return mailer })),
)

s.Client.POST("/api/v1/signup").JSON(...).Send().Status(201)

require.Len(t, mailer.Sent, 1)
require.Equal(t, "user@example.com", mailer.Sent[0].To)
```

The same pattern works for queue producers, feature-flag clients,
search indexes, and any other interface-shaped dependency. It runs
faster than HTTP mocking, gives type-safe call assertions, and isn't
sensitive to SDK upgrades.

> **Rule of thumb.** Mock at the *narrowest* interface your code
> depends on. If you only use `s3.PutObject`, define a one-method
> `Uploader` interface in your app and depend on that — your tests
> become one-line decorators instead of XML stub farms.

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
- `WithDBDecorator[T any](func(*gorm.DB) T)` — wrap the tx DB into your app's DB type (for simple wrappers built from `*gorm.DB` alone).
- `WithDBDecoratorFunc[T any](func(*gorm.DB, T) T)` — same idea, but the callback also receives the original instance fx built. Use this when your DB wrapper has unexported fields (logger, env, …) that you need to preserve.
- `WithEngineFrom[T any](func(T) *gin.Engine)` — tell the library where the engine lives.
- `WithAuthProvider(func(any) (name, value string))` — drives `Client.AsUser`.
- `WithSilentFxLogs()` — suppress fx lifecycle logging.

## Troubleshooting

**`fx.New failed: missing dependencies for function ... missing type: *gin.Engine`**
Your app exposes the engine inside a wrapper struct (e.g. `*infrastructure.Router`).
Add `gintest.WithEngineFrom(func(r *YourRouter) *gin.Engine { return r.Engine })`
to `NewSuite`.

**`missing type: *yourdb.Conn`**
The library injected `*gorm.DB` but your app expects a wrapper. Add
`gintest.WithDBDecorator(func(tx *gorm.DB) *yourdb.Conn { return &yourdb.Conn{DB: tx} })`.

**`viper: cannot read configuration: open .env: no such file`**
Your `framework.NewEnv` calls `log.Fatal` if `.env` is missing. Either
write a stub `.env` for tests (see [tests/setup.go](tests/setup.go) for
the pattern using `sync.Once`), or change `NewEnv` to tolerate a
missing file.

**Race detector flags writes to a package-global**
That's almost always an app-side singleton (loggers, env caches) being
written from parallel test setup. The library itself is race-clean;
fix the singleton with a `sync.Once`.

**Stubs from one test bleed into another**
`s.HTTP` patches a process-global transport. Don't use `t.Parallel()`
on HTTP-mock tests that target overlapping URLs. Or mock at the
interface level (Section C above) and skip HTTP entirely.

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
