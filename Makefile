COVERFILE := coverage.txt

lint:
	gofumpt -w ./..
	golangci-lint run --fix

run:
	go run ./.

clean:
	docker-compose down

build: clean
	docker-compose build

start: build
	docker-compose up -d

prune_volumes:
	docker volume prune --force

test: 
	go test -v -coverprofile=$(COVERFILE) ./...

go_check_deps:
	go list -u -m -json all

go_get_deps:
	go get -u ./...

init_deps:
	docker-compose up -d

prepare_deps:
	docker-compose down
	make prune_volumes
	make init_deps