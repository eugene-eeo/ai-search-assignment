build:
	cd anneal && go build && cd -

push:
	git push origin
	git push github
