default: bin/kaffeefee.tar.gz
all: clean bin/kaffeefee.tar.gz

.PHONY: clean
clean:
	rm -rf bin/

.PHONY: bin/kaffeefee
bin/kaffeefee:
	go get && CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o bin/kaffeefee .

container: bin/kaffeefee
	docker build --rm -t kaffeefee .

bin/kaffeefee.tar.gz: container
	docker save -o bin/kaffeefee.tar.gz kaffeefee
