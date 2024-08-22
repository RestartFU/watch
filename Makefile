install:
	./scripts/build.sh
push:
	git add .
	git commit -m "Update $( date +%D )"
	git push
view_page:
	librewolf https://github.com/restartfu/watch
