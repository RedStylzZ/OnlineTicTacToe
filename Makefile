build:
	@go build -o build/p2p main.go
	@GOOS=windows GOARCH=amd64 go build -o build/p2p.exe main.go

run: 
	@go build -o build/p2p main.go
	@./build/p2p