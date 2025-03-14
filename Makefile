BINARY=`git branch --show-current`

dev:
	@hivemind

test:
	go test ./...

assets:
	cd ui && \
	npm run build && \
	cd -

lint:
	golangci-lint run

linux:
	go build -o bin/server

windows:
	GOOS=windows GOARCH=amd64 go build -o bin/${BINARY}.exe

build_all: assets linux windows

prod: assets linux
	bin/server --dev --dir ./pb_data serve --http 0.0.0.0:8090
