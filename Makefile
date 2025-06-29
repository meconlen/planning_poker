# Planning Poker Docker Management

.PHONY: build run stop clean logs dev prod help

# Default target
help:
	@echo "Planning Poker Docker Commands:"
	@echo "  make build    - Build the Docker image"
	@echo "  make run      - Build and run the container"
	@echo "  make dev      - Run with docker-compose (development)"
	@echo "  make prod     - Run with docker-compose (production)"
	@echo "  make stop     - Stop all containers"
	@echo "  make clean    - Stop and remove containers/images"
	@echo "  make logs     - Show container logs"
	@echo "  make shell    - Open shell in running container"

# Build the Docker image
build:
	@echo "🐳 Building Planning Poker Docker image..."
	docker build -t planning-poker:latest .

# Build and run the container
run: build
	@echo "🚀 Starting Planning Poker container..."
	-docker stop planning-poker 2>/dev/null
	-docker rm planning-poker 2>/dev/null
	docker run -d \
		--name planning-poker \
		-p 8080:8080 \
		--restart unless-stopped \
		planning-poker:latest
	@echo "✅ Planning Poker is running at http://localhost:8080"

# Development environment with docker-compose
dev:
	@echo "🛠️ Starting development environment..."
	docker-compose up --build -d
	@echo "✅ Development environment running at http://localhost:8080"

# Production environment with docker-compose
prod:
	@echo "🚀 Starting production environment..."
	docker-compose -f docker-compose.prod.yml up --build -d
	@echo "✅ Production environment running at http://localhost:80"

# Stop all containers
stop:
	@echo "🛑 Stopping Planning Poker containers..."
	-docker stop planning-poker 2>/dev/null
	-docker-compose down 2>/dev/null
	-docker-compose -f docker-compose.prod.yml down 2>/dev/null

# Clean up containers and images
clean: stop
	@echo "🧹 Cleaning up containers and images..."
	-docker rm planning-poker 2>/dev/null
	-docker rmi planning-poker:latest 2>/dev/null
	-docker system prune -f

# Show container logs
logs:
	@echo "📋 Showing Planning Poker logs..."
	-docker logs -f planning-poker 2>/dev/null || docker-compose logs -f planning-poker

# Open shell in running container
shell:
	@echo "🐚 Opening shell in Planning Poker container..."
	docker exec -it planning-poker sh
