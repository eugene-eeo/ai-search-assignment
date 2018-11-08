build:
	cd anneal-2opt  && go build && cd -
	cd aco          && go build && cd -
	cd aco-parallel && go build && cd -

push:
	git push origin
	git push github

plot:
	python tools/hist.py 012 017 021 026 042 048 058 175 180 535
