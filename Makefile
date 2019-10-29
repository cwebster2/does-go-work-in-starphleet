GO_BUILD_ENV := CGO_ENABLED=0 GOOS=linux GOARCH=amd64

clean:
	echo "lol"

heroku: server.go
	go build .
