# 🚀 Go Microservice

[![Go Version](https://img.shields.io/badge/Go-1.21-blue.svg)](https://golang.org)
[![Docker](https://img.shields.io/badge/Docker-Ready-brightgreen.svg)](https://docker.com)
[![Kubernetes](https://img.shields.io/badge/Kubernetes-Compatible-326ce5.svg)](https://kubernetes.io)
[![API Documentation](https://img.shields.io/badge/API-Swagger-85EA2D.svg)](#api-documentation)

A production-ready Go microservice with advanced features including posts management, user authentication, social interactions, Redis caching, rate limiting, and comprehensive monitoring capabilities.

## 🌟 Features

- **RESTful API** - Clean REST endpoints with comprehensive CRUD operations
- **User Management** - Registration, authentication, and user profiles
- **Social Features** - Posts, comments, following/followers system
- **JWT Authentication** - Secure token-based authentication
- **Redis Caching** - High-performance caching layer
- **Rate Limiting** - Configurable request throttling
- **Database Migrations** - Version-controlled database schema
- **Email Notifications** - SendGrid integration for transactional emails
- **Health Checks** - Built-in health monitoring endpoints
- **Graceful Shutdown** - Production-ready server lifecycle management
- **Swagger Documentation** - Auto-generated API documentation
- **Docker & Kubernetes Ready** - Full containerization support

## 🛠 Tech Stack

- **Language**: Go 1.21
- **Database**: PostgreSQL
- **Cache**: Redis
- **Email Service**: SendGrid
- **Documentation**: Swagger/OpenAPI
- **Containerization**: Docker
- **Orchestration**: Kubernetes

## 🚀 Quick Start

### Prerequisites

- Go 1.21 or higher
- Docker (optional)
- PostgreSQL database
- Redis instance
- SendGrid API key

### Local Development

1. **Clone the repository**
   ```bash
   git clone git@github.com:vishal210893/Go_Microservice.git
   cd Go-Microservice

Install dependencies


go mod download
Set environment variables


export ADDR=:8000
export SENDGRID_API_KEY=your_sendgrid_api_key
export REDIS_ADDR=your_redis_address
export REDIS_PW=your_redis_password
export DB_ADDR=your_database_url
Run the application


go run cmd/api/main.go
The service will be available at http://localhost:8000
🐳 Docker Deployment
Build Docker Image

docker build -t go-microservice:latest .

Run with Docker

docker run -d \
--name go-microservice-container \
-p 8080:8000 \
-e ADDR=:8000 \
-e SENDGRID_API_KEY={$dummy_key} \
-e REDIS_ADDR=redis-10702.c264.ap-south-1-1.ec2.redns.redis-cloud.com:10702 \
-e REDIS_PW={pwd} \
vishal210893/go-microservice:1

Access the Application
Health Check: http://localhost:8080/v1/health
API Documentation: http://localhost:8080/swagger/index.html
Metrics: http://localhost:8080/debug/vars
Docker Management Commands


# View running containers
docker ps

# Check logs
docker logs go-microservice-container

# Stop container
docker stop go-microservice-container

# Remove container
docker rm go-microservice-container

☸️ Kubernetes Deployment
Deploy to Kubernetes
Apply ConfigMap


kubectl apply -f k8s/configmap.yaml
Apply Secrets


kubectl apply -f k8s/secret.yaml
Deploy the application


kubectl apply -f k8s/deployment.yaml
Create service
kubectl apply -f k8s/service.yaml
Kubernetes Files Structure

k8s/
├── configmap.yaml      # Environment configuration
├── secret.yaml         # Sensitive data (API keys, passwords)
├── deployment.yaml     # Application deployment
└── service.yaml        # Service and ingress configuration

Access via Kubernetes

# Check deployment status
kubectl get pods -l app=go-microservice

# Port forward for local access
kubectl port-forward service/go-microservice-service 8080:80

# View logs
kubectl logs -l app=go-microservice

**📁 Project Structure**

Go-Microservice/
├── cmd/
│   ├── api/               # API server implementation
│   └── migrate/           # Database migrations
├── internal/
│   ├── auth/              # Authentication logic
│   ├── db/                # Database connection
│   ├── env/               # Environment configuration
│   ├── mailer/            # Email service integration
│   ├── ratelimiter/       # Rate limiting implementation
│   └── repo/              # Repository pattern implementation
├── docs/                  # Swagger documentation
├── scripts/               # Build and deployment scripts
├── k8s/                   # Kubernetes manifests
├── Dockerfile            # Docker configuration
├── docker-compose.yml    # Local development setup
└── Makefile              # Build automation
Environment Variables
Variable
Description
Default
Required
ADDR
Server address
:8080
No
DB_ADDR
Database connection string
-
Yes
REDIS_ADDR
Redis server address
-
Yes
REDIS_PW
Redis password
-
Yes
SENDGRID_API_KEY
SendGrid API key for emails
-
Yes
JWT_SECRET
JWT signing secret
secret
Yes
RATE_LIMITER_ENABLED
Enable rate limiting
true
No
ENV
Environment (development/production)
development
No
📚 API Documentation
The API documentation is automatically generated using Swagger and available at:
Local: http://localhost:8000/swagger/index.html
Docker: http://localhost:8080/swagger/index.html
Key Endpoints
Health Check: GET /v1/health
User Registration: POST /v1/users
User Authentication: POST /v1/authentication/token
Posts: GET|POST|PUT|DELETE /v1/posts
Comments: GET|POST /v1/posts/{id}/comments
User Feed: GET /v1/users/{id}/feed
🧪 Testing

# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific test package
go test ./internal/auth

📊 Monitoring & Observability
The service includes built-in monitoring capabilities:


Health Endpoints: /v1/health for readiness and liveness probes
Metrics: /debug/vars for application metrics
Structured Logging: JSON-formatted logs for production
Graceful Shutdown: Handles SIGTERM/SIGINT signals properly
🤝 Contributing
Fork the repository
Create a feature branch (git checkout -b feature/amazing-feature)
Commit your changes (git commit -m 'Add some amazing feature')
Push to the branch (git push origin feature/amazing-feature)
Open a Pull Request
📄 License
This project is licensed under the MIT License - see the LICENSE file for details.


📞 Support
For support and questions:


Email: support@example.com
Issues: GitHub Issues
<hr></hr> Built with ❤️ using Go

This README provides comprehensive documentation covering all aspects of your microservice, including Docker and Kubernetes deployment instructions, project structure, API documentation, and proper formatting with badges and emojis for better visual appeal.
