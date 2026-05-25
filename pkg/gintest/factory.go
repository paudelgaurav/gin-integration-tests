package gintest

import (
	"sync/atomic"

	"gorm.io/gorm"
)

// Override mutates an in-progress model. Used by Factory.Build/Create to let
// callers customise individual fields per call.
type Override[T any] func(*T)

// Factory produces instances of T with sensible defaults. Inspired by
// factory_boy: define a builder once, customise per call.
//
//	var Project = gintest.NewFactory(func(seq int) models.Project {
//	    return models.Project{Name: fmt.Sprintf("Project %d", seq)}
//	})
//
//	p := Project.Create(db, gintest.Set(func(p *models.Project) {
//	    p.Endpoint = "https://example.com"
//	}))
type Factory[T any] struct {
	build func(seq int) T
	seq   atomic.Int64
}

// NewFactory creates a Factory whose builder receives an incrementing
// sequence number, useful for unique fields like names/emails.
func NewFactory[T any](build func(seq int) T) *Factory[T] {
	return &Factory[T]{build: build}
}

// Build produces an in-memory T with the given overrides applied. It does
// not persist anything.
func (f *Factory[T]) Build(overrides ...Override[T]) T {
	seq := int(f.seq.Add(1))
	v := f.build(seq)
	for _, o := range overrides {
		o(&v)
	}
	return v
}

// Create produces a T and persists it via the given DB. Fails the test if
// the insert errors.
func (f *Factory[T]) Create(t TestingT, db *gorm.DB, overrides ...Override[T]) T {
	t.Helper()
	v := f.Build(overrides...)
	if err := db.Create(&v).Error; err != nil {
		t.Fatalf("gintest: factory create: %v", err)
	}
	return v
}

// CreateN persists n instances and returns them.
func (f *Factory[T]) CreateN(t TestingT, db *gorm.DB, n int, overrides ...Override[T]) []T {
	t.Helper()
	out := make([]T, 0, n)
	for i := 0; i < n; i++ {
		out = append(out, f.Create(t, db, overrides...))
	}
	return out
}

// TestingT is the subset of *testing.T the factory uses. Defined as an
// interface so tests can pass mocks if they want.
type TestingT interface {
	Helper()
	Fatalf(format string, args ...any)
}
