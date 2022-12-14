CURRENT_DIR=$(shell pwd)

swag:
	swag init -g api/api.go -o api/docs

run:
	go run cmd/main.go

local-up:
	docker compose --env-file ./.env.docker up -d

proto-gen:
	rm -rf genproto
	./scripts/gen-proto.sh ${CURRENT_DIR}

pull-sub-module:
	git submodule update --init --recursive

update-sub-module:
	git submodule update --remote --merge

.PHONY: start