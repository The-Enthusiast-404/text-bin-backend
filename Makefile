# Load environment variables from .env file
ifneq (,$(wildcard .env))
    include .env
    export $(shell sed 's/=.*//' .env)
endif

# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

## confirm: ask for confirmation before running a command
.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]


# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

## run/api: run the Go API server
.PHONY: run/api
run/api:
	go run ./cmd/api -dsn=${DSN}

## db/psql: connect to the database using psql
.PHONY: db/psql
db/psql:
	psql ${DSN}

## db/migrations/new: create new migration files
.PHONY: db/migrations/new
db/migrations/new:
	@echo 'Creating migration files for ${name}...'
	migrate create -seq -ext=.sql -dir=./migrations ${name}

## db/migrations/up: run up migrations after confirmation
.PHONY: db/migrations/up
db/migrations/up: confirm
	@echo 'Running up migrations...'
	migrate -path ./migrations -database ${DSN} up


# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #

## audit: tidy dependencies and format, vet and test all code
.PHONY: audit
audit:
	@echo 'Tidying and verifying module dependencies...'
	go mod tidy
	go mod verify
	@echo 'Formatting code...'
	go fmt ./...
	@echo 'Vetting code...'
	go vet ./...
	staticcheck ./...
	@echo 'Running tests...'
	go test -race -vet=off ./...


# ==================================================================================== #
# BUILD
# ==================================================================================== #

current_time = $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
git_description = $(shell git describe --always --dirty --tags --long)

## build/api: build the cmd/api application
.PHONY: build/api
build/api:
	@echo 'Building cmd/api...'
	@echo 'Current time: $(current_time)'
	@echo 'Git description: $(git_description)'
	go build -ldflags="-s -X main.buildTime=$(current_time) -X main.version=$(git_description)" -o=./bin/api ./cmd/api
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -X main.buildTime=$(current_time) -X main.version=$(git_description)" -o=./bin/linux_amd64/api ./cmd/api


# ==================================================================================== #
# PRODUCTION
# ==================================================================================== #

production_host_ip = '142.93.222.185'

## production/connect: connect to the production server
.PHONY: production/connect
production/connect:
	ssh textbin@${production_host_ip}

## production/deploy/api: deploy the api to production
.PHONY: production/deploy/api
production/deploy/api:
	rsync -P ./bin/linux_amd64/api textbin@${production_host_ip}:~
	rsync -rP --delete ./migrations textbin@${production_host_ip}:~
	rsync -P ./remote/production/api.service textbin@${production_host_ip}:~
	rsync -P ./remote/production/Caddyfile textbin@${production_host_ip}:~
	ssh -t textbin@${production_host_ip} '\
		migrate -path ~/migrations -database $$DSN up \
        && sudo mv ~/api.service /etc/systemd/system/ \
        && sudo systemctl enable api \
        && sudo systemctl restart api \
        && sudo mv ~/Caddyfile /etc/caddy/ \
        && sudo systemctl reload caddy \
      '

## production/configure/caddyfile: configure the production Caddyfile
.PHONY: production/configure/caddyfile
production/configure/caddyfile:
	rsync -P ./remote/production/Caddyfile textbin@${production_host_ip}:~
	ssh -t textbin@${production_host_ip} '\
		sudo mv ~/Caddyfile /etc/caddy/ \
		&& sudo systemctl reload caddy \
	'
