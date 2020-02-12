VERSION=$(shell git rev-parse --short HEAD)

up-dev:
	docker-compose build --build-arg APP_VERSION=$(VERSION)
	docker-compose up 

run:
	docker build --build-arg APP_VERSION=$(VERSION) -t cabhelp_backend:dev .
	docker run --network host --publish 8088:8088 --env DATA_DIRECTORY=/go/src/cabhelp.ro/backend cabhelp_backend:dev