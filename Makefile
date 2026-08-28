include .env
export

DC             := docker compose
PG_SERVICE     := todo-app-postgres
MIGRATE_SERVICE := todo-app-postgres-migrate
FORWARD_SERVICE := port-forwarder

MIGRATE_CONTAINER := $(DC) run --rm $(MIGRATE_SERVICE)
DATABASE_URL       := postgresql://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(PG_SERVICE):5432/$(POSTGRES_DB)?sslmode=disable
MIGRATE            := $(MIGRATE_CONTAINER) -path /migrations -database "$(DATABASE_URL)"

ifneq ($(filter migrate-create,$(firstword $(MAKECMDGOALS))),)
  MIGRATE_NAME := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
  $(foreach w,$(MIGRATE_NAME),$(eval $(w):;@:))
endif

migrate-create:
	@if [ -z "$(MIGRATE_NAME)" ]; then \
		echo "Usage: make migrate-create <NAME>"; \
		exit 1; \
	fi
	@$(MIGRATE_CONTAINER) create -ext sql -dir /migrations -format 20060102150405 $(MIGRATE_NAME)

migrate-up:
	@$(MIGRATE) -verbose up

migrate-down:
	@$(MIGRATE) -verbose down 1

migrate-force-previous:
	@current=$$($(DC) exec -T $(PG_SERVICE) \
		psql "$(DATABASE_URL)" -qtAX \
		-c "SELECT version FROM schema_migrations LIMIT 1;" | tr -d '[:space:]'); \
	if [ -z "$$current" ]; then \
		echo "Error: could not read current migration version (is the DB up?)"; \
		exit 1; \
	fi; \
	previous=$$(find $(PROJECT_ROOT)/migrations -name '*.up.sql' \
		| sed -E 's|.*/([0-9]+)_.*\.up\.sql|\1|' \
		| sort -n \
		| awk -v current="$$current" '$$1 < current { prev=$$1 } END { print prev }'); \
	if [ -z "$$previous" ]; then \
		echo "Error: no migration older than $$current found"; \
		exit 1; \
	fi; \
	echo "Current migration: $$current"; \
	echo "Forcing to: $$previous"; \
	$(MIGRATE) force "$$previous"

env-up:
	@$(DC) up -d $(PG_SERVICE)

env-down:
	@$(DC) stop $(PG_SERVICE)

env-cleanup:
	@read -p "Erase the volume? There is a risk of data loss. [y/N]: " ans; \
	case "$$ans" in \
		[yY]) $(DC) rm -sf $(PG_SERVICE) $(FORWARD_SERVICE) && rm -rf $(PROJECT_ROOT)/out/pgdata ;; \
		*) echo "Aborted." ;; \
	esac

ifneq ($(filter env-port-forward,$(firstword $(MAKECMDGOALS))),)
  PORT_FORWARD := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
  $(foreach w,$(PORT_FORWARD),$(eval $(w):;@:))
endif

env-port-forward:
	@case "$(PORT_FORWARD)" in \
		'') echo "Usage: make env-port-forward <PORT>"; exit 1 ;; \
		*[!0-9]*) echo "Error: '$(PORT_FORWARD)' is not a valid port number"; exit 1 ;; \
	esac
	@if [ "$(PORT_FORWARD)" -lt 1 ] || [ "$(PORT_FORWARD)" -gt 65535 ]; then \
		echo "Error: port must be between 1 and 65535 (got $(PORT_FORWARD))"; \
		exit 1; \
	fi
	@PORT_FORWARD=$(PORT_FORWARD) $(DC) up -d $(FORWARD_SERVICE)

env-port-close:
	@$(DC) stop $(FORWARD_SERVICE)