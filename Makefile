build: clean
	go build -ldflags "-X main.debug true" -o server tcpecho.go
	go build -o client tcpecho.go

clean:
	-rm client server *.log
