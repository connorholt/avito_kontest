GO_VERSION := 1.21

APP_IMAGE := avito_app
DB_IMAGE := postgres:16

APP_CONTAINER := app
DB_CONTAINER := database


build:
	docker compose down
	docker build -t $(APP_IMAGE):1 ./
	docker compose up


