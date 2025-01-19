SERVER_APP_NAME=bin/in-mem-kvdb-server

build-server:
	go build -o ${SERVER_APP_NAME} cmd/server/main.go

run-server-with-config: build-server
	CONFIG_FILE_NAME=database_config.yml ./${SERVER_APP_NAME}
