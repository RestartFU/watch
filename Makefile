install:
	sudo ./scripts/build.sh
push:
	git add .
	git commit -m "Update $( date +%D )"
	git push
view_page:
	librewolf https://github.com/restartfu/watch
test_go:
	sudo ./scripts/build.sh
	cd examples/go && watch
test:
	python scripts/test_all.py
