.PHONY: install build debugBuild clean deploy

install:
	npm i

build:
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/hello hello/main.go
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/world world/main.go
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/scraper scraper/main.go

debugBuild:
	env GOARCH=amd64 GOOS=linux go build -v -gcflags='all=-N -l' -ldflags="-s -w" -o bin/hello hello/main.go
	env GOARCH=amd64 GOOS=linux go build -v -gcflags='all=-N -l' -ldflags="-s -w" -o bin/world world/main.go
	env GOARCH=amd64 GOOS=linux go build -v -gcflags='all=-N -l' -ldflags="-s -w" -o bin/scraper scraper/main.go

debugDeploy: clean install debugBuild
	env SLS_DEBUG=* npx sls offline --useDocker

clean:
	rm -rf ./node_modules
	rm -rf ./bin

deploy: clean install build
	npx sls deploy --verbose
