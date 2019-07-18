#/bin/bash
NAME=name/reminder
BOT_NAME=Od15x4Bot
DSN=postgres://bot:@bot-db:5432/bot?sslmode=disable
NETWORK=15x4
run: up migrate logs
startup: create-network run
create-network:
	-@docker network create -d bridge $(NETWORK)
logs: 
	docker-compose logs -f
up: 
	docker-compose up --build --force-recreate -d 
down:
	docker-compose stop $(call name,'')	
exec:	
	docker-compose exec $(call name,'')	bash
test:
	-@docker-compose run \
		--rm \
		bot \
		go test $(if $(params),$(params),$(shell echo '-cover -race -v ./...'))
shutdown:
	# stop containers, remove volumes and containers for services not defined in the compose file
	docker-compose down --remove-orphans -v

# migration commands 
migrate: 
	docker run -v $(PWD)/postgresql/migrations:/migrations --network=$(NETWORK)  migrate/migrate -path=/migrations/ -database $(DSN) up

migrate-create: 
	docker run -v $(PWD)/postgresql/migrations:/migrations --network=$(NETWORK)  migrate/migrate create -dir=/migrations/ -ext=.sql $(name)

migrate-down:
	docker run  -v $(PWD)/postgresql/migrations:/migrations --network=$(NETWORK)  migrate/migrate -path=/migrations/ -database $(DSN) down 1

migrate-version: 
	docker run -v $(PWD)/postgresql/migrations:/migrations --network=$(NETWORK)  migrate/migrate -path=/migrations/ -database $(DSN) version

migrate-fix:
	docker run  -v $(PWD)/postgresql/migrations:/migrations --network=$(NETWORK)  migrate/migrate -path=/migrations/ -database $(DSN) force $(migration_version)
