build:
	cd anneal && go build && cd -
	cd aco    && go build && cd -

push:
	git push origin
	git push github

plot:
	python tools/hist.py 012 017 021 026 042 048 058 175 180 535
