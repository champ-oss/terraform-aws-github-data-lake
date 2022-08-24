pytest:
	cd github_event_receiver_lambda && coverage run -m pytest

flake8:
	cd github_event_receiver_lambda && flake8 . --count --max-complexity=12 --max-line-length=127 --statistics --exclude venv