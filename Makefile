
tidy:
	go mod tidy

build.server:
	cd cmd/server && go build -o server .

run.server: build.server
	cd cmd/server && ./server