all:
	go build -o stockticker

install:
	go install

clean:
	rm -f $(GOPATH)/bin/stockticker
