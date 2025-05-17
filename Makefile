
run:
	docker-compose up --build

test:
	go test ./... -v

lint:
	golangci-lint run

build:
	go build -o bin/app main.go

clean:
	rm -rf bin

docker-down:
	docker-compose down
