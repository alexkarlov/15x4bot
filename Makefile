#/bin/bash
NAME=name/reminder
BOT_NAME=Od15x4Bot
run: up logs
logs: 
	docker-compose logs -f
up: 
	docker-compose up --build --force-recreate -d 
down:
	docker-compose down --remove-orphans -v $(call name,'')	
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
