PWD=$$(pwd)

build:
	yarn install
	yarn build
	mage -v

rm:
	docker stop grafana
	docker rm grafana

run:
	docker run -d -p 3000:3000 -v ${PWD}:/var/lib/grafana/plugins --name=grafana grafana/grafana:7.0.0

sign:
	export GRAFANA_API_KEY=eyJrIjoiYmMxZjFiMGQ5MzdmMzI4NDZlY2E4ZTdmMzViYmU5MDFhOWI3YzhkOSIsIm4iOiJ1Y2xvdWQtZ3JhZmFuYSIsImlkIjo1MjMzMDd9 \
	&& npx @grafana/toolkit plugin:sign --rootUrls http://127.0.0.1:3000

restart:
	docker restart grafana

all: lint build sign rm run

dev: lint build

lint:
	npx prettier --write ./src/**