
tidy:
	go mod tidy

build.server:
	cd cmd/server && go build -o ../out/server .

run.server: build.server
	cd cmd/out && ./server