.PHONY: clean cleanServerless installDebugTools installServerless build buildDocker buildDebugDocker runDebug runDebugDocker deploy

clean:
	rm -rf ./bin

cleanServerless:
	rm ./node_modules

installServerless:
	npm i

installDebugTools:
	wget https://github.com/aws/aws-lambda-runtime-interface-emulator/releases/latest/download/aws-lambda-rie
	mkdir -p .aws-lambda-rie 
	mv aws-lambda-rie .aws-lambda-rie/aws-lambda-rie
	chmod +x .aws-lambda-rie/aws-lambda-rie

build:
	GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/scraper scraper/main.go

debug:
	env GOARCH=amd64 GOOS=linux go build -v -gcflags='all=-N -l' -ldflags="-s -w" -o bin/scraperDebug scraper/main.go

runDebug: debug
	DEBUG=true ./.aws-lambda-rie/aws-lambda-rie ./bin/scraperDebug

buildDocker:
	docker build -t watcher-local-build .

buildDebugDocker:
	docker build --build-arg="DEBUG=true" -t watcher-local-build .

runDebugDocker: buildDebugDocker
	touch page-docker.html
	docker run --platform linux/amd64 -v ./.aws-lambda-rie:/aws-lambda -v ./page-docker.html:/task/page.html -p 0.0.0.0:9000:8080 --entrypoint /aws-lambda/aws-lambda-rie watcher-local-build /task/app

deploy: install
	npx sls deploy --verbose
