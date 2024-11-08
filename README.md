## Tests are included in the `tests` folder.

This repository is an example repo, not a package.  
The provided sample can be found here:  
https://github.com/paudelgaurav/gin-integration-tests/blob/master/tests/project_endpoint_test.go

### TODO
- [ ] Set up a test database (the current example uses SQLite; needs refactoring to MySQL)
- [ ] Implement authentication handling in the test runner
- [ ] Handle body data that may include actual IDs from other related tables
- [ ] Assert response body contents along with status codes
- [ ] Add assertion functions to verify data creation
- [ ] Enable the option to delete or preserve the test database
- [ ] Replace `.env` values specifically for tests
- [ ] Mock infrastructure dependencies such as AWS or GCP services
