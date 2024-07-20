up:
	docker compose up -d

down:
	docker compose down

rebuild:
	docker build -t todo-list .

.PHONY: up down restart rebuild