OTEL_RESOURCE_ATTRIBUTES="service.name=dice,service.version=0.1.0"


.PHONY: build
build:
	sam build --no-cached


.PHONY: run
run: build
	sam local invoke --event "./events/event.json"