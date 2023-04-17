build:
	go build -o bin/mysql-ps -ldflags="-s -w -extldflags=-static" main.go

build-all: build

tar: build
	tar --directory=. --transform='s|bin||' -czvf build.tgz bin/mysql-ps .env.example

clean:
	rm -vf bin/mysql-ps *.tgz