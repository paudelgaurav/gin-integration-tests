## Tests are included in the `tests` folder.

This repository is an example repo, not a package.  
The provided sample can be found here:  
https://github.com/paudelgaurav/gin-integration-tests/blob/master/tests/project_endpoint_test.go

### New Features ✨

#### Response Body Assertions
The test framework now supports comprehensive response body validation alongside status code assertions:

- **Exact Response Body Matching**: Validate entire JSON structure
- **Partial Content Matching**: Validate specific fields and nested objects
- **Custom Assertion Functions**: Implement complex validation logic

See [Response Body Assertions Documentation](./docs/RESPONSE_BODY_ASSERTIONS.md) for detailed usage examples.

### TODO
- [ ] Set up a test database (the current example uses SQLite; user should be able to add mysql and postgres)
- [ ] Implement authentication handling in the test runner
- [X] Handle body data that may include actual IDs from other related tables (currently handled with prepareBodyFunc
- [X] Assert response body contents along with status codes
- [ ] Add assertion functions to verify data creation
- [ ] Enable the option to delete or preserve the test database
- [ ] Replace `.env` values specifically for tests
- [ ] Mock infrastructure dependencies such as AWS or GCP services
