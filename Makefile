APP_NAME=bins/in-mem-kvdb


build-app:
	go build -o $(APP_NAME) cmd/main.go

run-app: build-app
	CONFIG_FILE_NAME=config.yml ./$(APP_NAME) kvdb-cli
