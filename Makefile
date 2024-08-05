.PHONY: build

build:
	sam build


.PHONY: run
run: build
	sam local invoke --event "./events/event.json"