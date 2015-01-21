install:
	go install

remove:
	rm -f $(GOPATH)/bin/stockticker

clean:
	rm -f stockticker
