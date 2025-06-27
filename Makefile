run:
	go run cmd/server/main.go
run-docker:
	docker-compose up -d
stop-docker:
	docker-compose down