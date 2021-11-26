
.PHONY: test leapseconds.txt

leapseconds.txt:
	curl https://www.ietf.org/timezones/data/leap-seconds.list | grep "^[^#]" > leapseconds.txt
	
test:
	go test  ./...
