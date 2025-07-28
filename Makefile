u:
	docker-compose up -d

dw:
	docker-compose -f docker-compose.yml down --remove-orphans

ps:
	docker-compose ps