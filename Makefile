help :
	@echo "Help information, please run specific target:"
	@IFS=$$'\n' ; \
	help_lines=(`fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//'`); \
	for help_line in $${help_lines[@]}; do \
		IFS=$$'#' ; \
		help_split=($$help_line) ; \
		help_command=`echo $${help_split[0]} | sed -e 's/^ *//' -e 's/ *$$//'` ; \
		help_info=`echo $${help_split[2]} | sed -e 's/^ *//' -e 's/ *$$//'` ; \
		printf " - %-20s %s\n" $$help_command $$help_info ; \
	done

api-doc: ## Generate API DOC
	@swag init -g cmd/api/main.go

test: ## Run the tests
	@go test -v ./...

install: ## Setting up API
	@docker compose down -v
	@docker compose build
	@docker compose up -d