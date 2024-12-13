.PHONY: venv install test workflow

venv:
	virtualenv venv

install:
	venv/bin/pip install -r blog/requirements.txt && venv/bin/pip install pytest

test:
	venv/bin/python -m pytest blog/test.py

workflow: venv install test
cleanup:
	rm -rf venv && rm test.db
