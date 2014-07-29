all: build

build:
	mkdir -p bin
	go build -o bin/server  github.com/oikomi/gortmpserver/server
	go build -o bin/client  github.com/oikomi/gortmpserver/client

clean:
	rm -rf bin
