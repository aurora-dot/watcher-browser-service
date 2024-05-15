.PHONY: clean install build buildDocker runDebug runDebugDocker getDebugTools deploy

clean:
	rm -rf ./node_modules

install:
	npm i

getDebugTools:
	wget https://github.com/aws/aws-lambda-runtime-interface-emulator/releases/latest/download/aws-lambda-rie
	mkdir -p .aws-lambda-rie 
	mv aws-lambda-rie .aws-lambda-rie/aws-lambda-rie
	chmod +x .aws-lambda-rie/aws-lambda-rie

build:
	GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/scraper scraper/main.go

debug:
	env GOARCH=amd64 GOOS=linux go build -v -gcflags='all=-N -l' -ldflags="-s -w" -o bin/scraperDebug scraper/main.go

runDebug: debug
	./.aws-lambda-rie/aws-lambda-rie ./bin/scraperDebug

buildDocker:
	docker build -t watcher-local-build .

runDebugDocker: buildDocker
	docker run --platform linux/amd64 -v ./.aws-lambda-rie:/aws-lambda -p 9000:8080 --entrypoint /aws-lambda/aws-lambda-rie watcher-local-build /task/app

deploy: install
	npx sls deploy --verbose
