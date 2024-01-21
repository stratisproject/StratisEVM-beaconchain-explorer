start-infra:
	docker compose -f docker-compose.infra.yml up -d --build

stop-infra:
	docker compose -f docker-compose.infra.yml down

start-explorer:
	docker compose -f docker-compose.yml up -d --build

stop-explorer:
	docker compose -f docker-compose.yml down

explorer-logs:
	docker compose -f docker-compose.yml logs -f
