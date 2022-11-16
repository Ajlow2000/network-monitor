build:
	docker build --rm -t network-monitor:test .

run:
	docker run --rm network-monitor:test