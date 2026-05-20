package gintest

// AssertCount fails the test unless the table backing model has exactly n rows
// matching the optional where clause. The where clause uses GORM's variadic
// query syntax: AssertCount(&Project{}, 2, "name LIKE ?", "Foo%").
func (s *Suite) AssertCount(model any, n int64, where ...any) {
	s.T.Helper()
	q := s.DB.Model(model)
	if len(where) > 0 {
		q = q.Where(where[0], where[1:]...)
	}
	var got int64
	if err := q.Count(&got).Error; err != nil {
		s.T.Fatalf("gintest: AssertCount: %v", err)
	}
	if got != n {
		s.T.Fatalf("gintest: AssertCount expected %d, got %d", n, got)
	}
}

// AssertExists fails the test unless at least one row matches.
func (s *Suite) AssertExists(model any, where ...any) {
	s.T.Helper()
	q := s.DB.Model(model)
	if len(where) > 0 {
		q = q.Where(where[0], where[1:]...)
	}
	var got int64
	if err := q.Count(&got).Error; err != nil {
		s.T.Fatalf("gintest: AssertExists: %v", err)
	}
	if got == 0 {
		s.T.Fatalf("gintest: AssertExists: no rows matched")
	}
}

// AssertNotExists fails the test if any row matches.
func (s *Suite) AssertNotExists(model any, where ...any) {
	s.T.Helper()
	q := s.DB.Model(model)
	if len(where) > 0 {
		q = q.Where(where[0], where[1:]...)
	}
	var got int64
	if err := q.Count(&got).Error; err != nil {
		s.T.Fatalf("gintest: AssertNotExists: %v", err)
	}
	if got != 0 {
		s.T.Fatalf("gintest: AssertNotExists: expected 0 rows, found %d", got)
	}
}
