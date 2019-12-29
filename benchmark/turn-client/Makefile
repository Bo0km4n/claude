.PHONY: push
push:
	git add .
	git commit -m "Update turn clients"
	git push origin master
pull:
	git pull origin master

.PHONY:bin
bin:
	GOOS=linux GOARCH=amd64 go build -o ./bin/reader.bin ./reader/main.go
	GOOS=linux GOARCH=amd64 go build -o ./bin/writer.bin ./writer/main.go
	GOOS=linux GOARCH=amd64 go build -o ./bin/ping.bin ./ping/main.go
clean:
	rm ./bin/*.bin