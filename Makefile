build:
	cd server
	go build -o ../out/server.exe ./cmd/
	cd ..
	cd client
	go build -o ../out/client.exe ./cmd/
	cd ..