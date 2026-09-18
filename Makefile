include .env
export

DC              := docker compose
PG_SERVICE      := todo-app-postgres
MIGRATE_SERVICE := todo-app-postgres-migrate
FORWARD_SERVICE := port-forwarder

PASSWORD_ENCODER    := $(PROJECT_ROOT)/scripts/encode_db_password.sh

DATABASE_PASSWORD   := $(shell $(PASSWORD_ENCODER) $(POSTGRES_PASSWORD_FILE))
DATABASE_URL        := postgresql://$(POSTGRES_USER):$(DATABASE_PASSWORD)@$(PG_SERVICE):5432/$(POSTGRES_DB)?sslmode=disable

MIGRATE             := $(DC) run --rm $(MIGRATE_SERVICE)

.PHONY: guard-% check-postgres-env check-project-root-env check-postgres-password \
        migrate-create migrate-up migrate-down migrate-force-previous \
        env-up env-down env-cleanup env-port-forward env-port-close

guard-%:
	@if [ -z '$($(*))' ]; then \
		echo "Error: required variable $* is not set"; \
		exit 1; \
	fi

check-postgres-env: guard-POSTGRES_USER guard-POSTGRES_PASSWORD_FILE guard-POSTGRES_DB

check-project-root-env: guard-PROJECT_ROOT

#ifneq ($(filter migrate-create,$(firstword $(MAKECMDGOALS))),)
#  MIGRATE_NAME := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
#  $(foreach w,$(MIGRATE_NAME),$(eval $(w):;@:))
#endif
#migrate-create:
#	@if [ -z "$(MIGRATE_NAME)" ]; then \
#		echo "Usage: make migrate-create <NAME>"; \
#		exit 1; \
#	fi
#	@$(MIGRATE_CONTAINER) create -ext sql -dir /migrations -format 20060102150405 $(MIGRATE_NAME)

check-postgres-password:
	@$(PASSWORD_ENCODER) $(POSTGRES_PASSWORD_FILE) >/dev/null

migrate-create: check-project-root-env
	@if [ -z "$(seq)" ]; then \
		echo "Usage: make migrate-create seq=<NAME>"; \
		exit 1; \
	fi
	@$(MIGRATE) create -ext sql -dir /migrations -format 20060102150405 $(seq)

migrate-up: check-project-root-env check-postgres-env check-postgres-password
	@$(MIGRATE) -verbose up

migrate-down: check-project-root-env check-postgres-env check-postgres-password
	@$(MIGRATE) -verbose down 1

migrate-force-previous: check-project-root-env check-postgres-env check-postgres-password
	@row=$$($(DC) exec -T $(PG_SERVICE) \
		psql "$(DATABASE_URL)" -qtAX \
		-c "SELECT version, dirty FROM schema_migrations LIMIT 1;" | tr -d '[:space:]'); \
	current=$$(echo "$$row" | cut -d'|' -f1); \
	dirty=$$(echo "$$row" | cut -d'|' -f2); \
	if [ -z "$$current" ]; then \
		echo "Error: could not read current migration version (is the DB up?)"; \
		exit 1; \
	fi; \
	if [ "$$dirty" != "t" ]; then \
		echo "Error: database is not dirty (current version: $$current)\nNothing to force"; \
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
	echo "Current migration: $$current (dirty)"; \
	echo "Forcing to: $$previous"; \
	$(MIGRATE) force "$$previous"

env-up: check-project-root-env check-postgres-env check-postgres-password
	@$(DC) up -d $(PG_SERVICE)

env-down:
	@$(DC) down $(PG_SERVICE) $(FORWARD_SERVICE)

env-cleanup: check-project-root-env
	@read -p "Erase the volume? There is a risk of data loss [y/N]: " ans; \
	case "$$ans" in \
		[yY]) $(DC) rm -sf $(PG_SERVICE) $(FORWARD_SERVICE) && rm -rf $(PROJECT_ROOT)/out/pgdata ;; \
		*) echo "Aborted" ;; \
	esac

ifneq ($(filter env-port-forward,$(firstword $(MAKECMDGOALS))),)
  PORT_FORWARD := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
  $(foreach w,$(PORT_FORWARD),$(eval $(w):;@:))
endif

env-port-forward: check-postgres-env
	@case "$(PORT_FORWARD)" in \
		'') echo "Usage: make env-port-forward <PORT>"; exit 1 ;; \
		*[!0-9]*) echo "Error: '$(PORT_FORWARD)' is not a valid port number"; exit 1 ;; \
	esac
	@if [ "$(PORT_FORWARD)" -lt 1 ] || [ "$(PORT_FORWARD)" -gt 65535 ]; then \
		echo "Error: port must be between 1 and 65535 (got $(PORT_FORWARD))"; \
		exit 1; \
	fi
	@PORT_FORWARD=$(PORT_FORWARD) $(DC) up -d --wait $(FORWARD_SERVICE) || { \
		echo "Error: port-forwarder failed to become healthy — check 'docker compose logs $(FORWARD_SERVICE)'"; \
		exit 1; \
	}

env-port-close:
	@$(DC) stop $(FORWARD_SERVICE)
