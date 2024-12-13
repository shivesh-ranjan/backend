install:
	python -m pip install -r blog/requirements.txt
install-pytest:
	python -m pip install pytest
test:
	python -m pytest blog/test.py
