#/bin/bash
NAME=name/reminder
BOT_NAME=Od15x4Bot
run: up logs
logs: 
	docker-compose logs -f
up: 
	docker-compose up --build --force-recreate -d 
shutdown:
	docker-compose down --remove-orphans -v $(call name,'')	
exec:	
	docker-compose exec $(call name,'')	bash