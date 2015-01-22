install:
	go install

remove:
	go clean -n -i -x
	rm -f $(GOPATH)/bin/stockticker

clean:
	rm -f stockticker
