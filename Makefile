.PHONY: clean install build push deploy

clean:
	rm -rf ./node_modules

install:
	npm i

deploy: install
	npx sls deploy --verbose
