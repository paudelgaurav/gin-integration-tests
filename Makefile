include .env
export

lint-setup:
	python3 -m ensurepip --upgrade
	sudo pip3 install pre-commit
	pre-commit install
	pre-commit autoupdate

.PHONY: lint-setup