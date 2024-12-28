
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
	@swag init -g cmd/postgresql/main.go

test: ## Run the tests
	@go test -v ./...

make install:  ## Build and up project
	@docker compose up -d --build

#migrate: ## Roll out migrations (only once)
#	migrate -path ./schema -database 'postgres://postgres:qwerty@0.0.0.0:5438/postgres?sslmode=disable' up