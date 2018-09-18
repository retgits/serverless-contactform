# Description: Makefile for AWS Lambda functions 
# Author: Leon Stigter <lstigter@gmail.com>
# Last Updated: 2018-09-18

.PHONY: deps clean build deploy test-lambda

# Variables
FUNCTION=contactform
S3=retgits-apps

deps:
	go get -u=patch ./...

clean: 
	rm -rf ./bin
	
build:
	GOOS=linux GOARCH=amd64 go build -o ./bin/$(FUNCTION)-lambda *.go

test-lambda: clean build
	sam local invoke $(FUNCTION) -e ./test/event.json

deploy: clean build
	sam package --template-file template.yaml --output-template-file packaged.yaml --s3-bucket $(S3)
	sam deploy --template-file packaged.yaml --stack-name $(FUNCTION)-lambda --capabilities CAPABILITY_IAM