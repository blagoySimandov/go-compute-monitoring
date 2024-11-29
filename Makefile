start:
	@echo "Starting server..."
	@go run cmd/server/main.go

run-all:
	@docker-compose up -d prometheus grafana matrix-compute

local-start:
	@echo "Starting server locally..."
	@go run cmd/server/main.go

local-monitoring:
	@echo "Starting Prometheus and Grafana containers..."
	@docker-compose up prometheus grafana

# Run application locally with containerized monitoring
local-all:
	@docker-compose up -d prometheus grafana
	@make local-start

stop-monitoring:
	@echo "Stopping monitoring containers..."
	@docker-compose stop prometheus grafana

clean:
	@echo "Cleaning up local data..."
	@rm -rf ./data
