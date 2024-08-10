include .env
export $(shell sed 's/=.*//' .env)

# CURRENT_DIR va DB_URL ni yaratish
CURRENT_DIR := $(shell pwd)
DB_URL := postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable

# Migratsiya yo'li va mig komandalari
MIGRATE_PATH := database/migrations
MIGRATE_CMD := migrate -path $(MIGRATE_PATH) -database '$(DB_URL)' -verbose

.PHONY: gen-proto mig-up mig-down mig-force mig-create

gen-proto:
	./scripts/gen-proto.sh $(CURRENT_DIR)

mig-up:
	$(MIGRATE_CMD) up

mig-down:
	$(MIGRATE_CMD) down

mig-force:
	$(MIGRATE_CMD) force 1

mig-create:
	migrate create -ext sql -dir $(MIGRATE_PATH) -seq personal_finance_tracker_table

# Swaggerni generate qilish
swag-init:
	~/go/bin/swag init -g ./api/routes.go -o api/docs