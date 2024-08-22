install:
	sudo ./scripts/build.sh
push:
	./scripts/push.sh
view_page:
	librewolf https://github.com/restartfu/watch
test_go:
	sudo ./scripts/build.sh
	cd examples/go && watch
test:
	rm -rf bin/*
	python scripts/test_all.py
