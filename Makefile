.PHONY: help build up down logs clean restart

help:
	@echo "miniRAG Server - Available commands:"
	@echo "  make build      - Build Docker image"
	@echo "  make up         - Start all services"
	@echo "  make down       - Stop all services"
	@echo "  make logs       - View logs"
	@echo "  make restart    - Restart services"
	@echo "  make clean      - Remove containers and volumes"
	@echo "  make test       - Run API tests"

build:
	docker-compose build

up:
	docker-compose up -d
	@echo "✓ Services started"
	@echo "  App UI:   http://localhost:80"
	@echo "  API:      http://localhost:8080"
	@echo "  Weaviate: http://localhost:8081"

down:
	docker-compose down

logs:
	docker-compose logs -f minirag

restart:
	docker-compose restart

clean:
	docker-compose down -v
	@echo "✓ All containers and volumes removed"

test:
	@echo "Adding test documents..."
	curl -X POST http://localhost:8080/adddocument \
	  -H "Content-Type: application/json" \
	  -d '{"documents": [{"text": "Go is a compiled language"}, {"text": "Python is interpreted"}]}'
	@echo "\n\nQuerying..."
	curl -X POST http://localhost:8080/queryprompt \
	  -H "Content-Type: application/json" \
	  -d '{"content": "What is Go?"}'

status:
	docker-compose ps

shell:
	docker-compose exec minirag /bin/sh
