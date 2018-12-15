.PHONY: 

#--- Help ---
help:
	@echo 
	@echo Makefile targets
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' Makefile | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
	@echo

#--- deps ---
deps: ## Get and update all dependencies
	go get -u=patch ./...

#--- test ---
test: ## Run go test
	. .env \
	&& go test ./...

#--- update-secrets ---
update-secrets: ## Update the secrets
	. .env && now secrets rm recaptcha-secret -y && now secrets add recaptcha-secret $$RECAPTCHA_SECRET
	. .env && now secrets rm email-address && now secrets add email-address $$EMAIL_ADDRESS
	. .env && now secrets rm email-password && now secrets add email-password $$EMAIL_PASSWORD

#--- deploy ---
deploy: ## deploy the app to Zeit
	now