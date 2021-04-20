
tidy:
	go mod tidy

build.server:
	cd cmd/server && go build -o ../out/server .

run.server: build.server
	cd cmd/out && ./server


build.client:
	cd cmd/client && go build -o ../out/client .

run.client: build.client
	cd cmd/out && ./client

clean: tidy
	rm -rf cmd/out