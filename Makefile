
up:
	docker-compose up --build

down:
	docker-compose down

test:
	docker-compose exec app go test ./...
