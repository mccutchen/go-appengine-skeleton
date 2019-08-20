PROJECT_ROOT := $(dir $(realpath $(firstword $(MAKEFILE_LIST))))
PROJECT := $(lastword $(subst /, , $(PROJECT_ROOT)))

APP_PORT := 9999

DS_PATH := $(PROJECT_ROOT)/.datastore
DS_PORT := 9111
DS_HOST := 0.0.0.0:$(DS_PORT)
DS_CONTAINER_NAME := appengine-datastore-$(PROJECT)

BIN_PATH := dist/$(PROJECT)

build: *.go go.*
	mkdir -p dist
	go build -o $(BIN_PATH)

test:
	go test -v -timeout 30s

integrationtest:
	go test -v -timeout 5s -tags=integration

run: build
	DATASTORE_EMULATOR_HOST=$(DS_HOST) \
	DATASTORE_PROJECT_ID=$(PROJECT) \
	$(BIN_PATH)

devdb:
	# https://cloud.google.com/datastore/docs/tools/datastore-emulator
	mkdir -p $(DS_PATH)
	docker run \
		--rm \
		--name $(DS_CONTAINER_NAME) \
		-p $(DS_PORT):$(DS_PORT) \
		-v $(DS_PATH):/data \
		-v $(HOME)/.config/gcloud:/root/.config/gcloud \
		google/cloud-sdk \
		gcloud beta emulators datastore start --project=$(PROJECT) --host-port $(DS_HOST) --data-dir=/data

ngrok:
	ngrok http -subdomain=$(PROJECT) $(APP_PORT)
