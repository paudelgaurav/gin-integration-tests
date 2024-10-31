include .env
export

MIGRATE=atlas migrate

migrate-status:
	$(MIGRATE) status --url "sqlite://data.db"

migrate-diff:
	${MIGRATE} diff --env gorm

migrate-apply:
	$(MIGRATE) apply --url "sqlite://data.db"

apply-test-migrations:
	$(MIGRATE) apply --url "sqlite://tests/data_test.db"

migrate-down:
	$(MIGRATE) down --url "sqlite://data.db"

migrate-hash:
	$(MIGRATE) hash

lint-setup:
	python3 -m ensurepip --upgrade
	sudo pip3 install pre-commit
	pre-commit install
	pre-commit autoupdate

.PHONY: migrate-status migrate-diff migrate-apply migrate-down migrate-hash lint-setup apply-test-migrations