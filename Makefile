install:
	./script/build.sh
push:
	git add .
	git commit -m "Update $( date +%D )"
	git push
