build:
	cd anneal-2opt  && go build && cd -
	cd aco          && go build && cd -
	cd aco-parallel && go build && cd -

push:
	git push origin
	git push github
