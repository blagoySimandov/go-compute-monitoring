install-prometheus:
	@echo "Downloading and setting up Prometheus..."
	@wget https://github.com/prometheus/prometheus/releases/download/v2.49.1/prometheus-2.49.1.linux-amd64.tar.gz
	@tar xvfz prometheus-*.tar.gz
	@mv prometheus-*-linux-amd64 prometheus
	@rm prometheus-*.tar.gz
	@echo "Prometheus installed successfully"

install-go:
	@echo "Installing Go..."
	@wget https://go.dev/dl/go1.22.1.linux-amd64.tar.gz
	@sudo rm -rf /usr/local/go
	@sudo tar -C /usr/local -xzf go1.22.1.linux-amd64.tar.gz
	@rm go1.22.1.linux-amd64.tar.gz
	@echo "export PATH=$$PATH:/usr/local/go/bin" >> ~/.bashrc
	@source ~/.bashrc
	@echo "Go installed successfully"

setup: install-go install-prometheus
	@echo "Installing project dependencies..."
	@go mod download
	@echo "Setup complete!"

run-prometheus:
	@echo "Starting Prometheus..."
	@./prometheus/prometheus --config.file=prometheus.yml

start:
	@echo "Starting server..."
	@go run cmd/server/main.go

run-all:
	@make run-prometheus & make start
