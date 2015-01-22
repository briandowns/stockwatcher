install:
	go install

remove:
	go clean -n -i -x
	rm -f $(GOPATH)/bin/stockwatcher

clean:
	rm -f stockwatcher
