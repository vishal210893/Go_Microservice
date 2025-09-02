# 🚀 Go Microservice

<div align="center">

[![Go Version](https://img.shields.io/badge/Go-1.21-blue.svg)](https://golang.org)
[![Docker](https://img.shields.io/badge/Docker-Ready-brightgreen.svg)](https://docker.com)
[![Kubernetes](https://img.shields.io/badge/Kubernetes-Compatible-326ce5.svg)](https://kubernetes.io)
[![API Documentation](https://img.shields.io/badge/API-Swagger-85EA2D.svg)](#-api-documentation)
[![License](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

*A production-ready Go microservice with advanced features including social interactions, caching, rate limiting, and comprehensive monitoring capabilities.*

[Getting Started](#-quick-start) • [API Documentation](#-api-documentation) • [Docker Deployment](#-docker-deployment) • [Kubernetes](#️-kubernetes-deployment)

</div>

---

## 🛠 Tech Stack

| Category | Technology |
|----------|------------|
| **Language** | Go 1.21+ |
| **Database** | PostgreSQL with connection pooling |
| **Cache** | Redis with clustering support |
| **Authentication** | JWT with HMAC-SHA256 |
| **Email Service** | SendGrid API |
| **Documentation** | Swagger/OpenAPI 3.0 |
| **Containerization** | Docker with multi-stage builds |
| **Orchestration** | Kubernetes with Helm charts |
| **Monitoring** | Built-in metrics and health checks |

---

## 🚀 Quick Start

### Prerequisites

Ensure you have the following installed on your system:

- **Go 1.21+** - [Download here](https://golang.org/dl/)
- **PostgreSQL 13+** - [Installation guide](https://www.postgresql.org/download/)
- **Redis 6+** - [Installation guide](https://redis.io/download)
- **Docker** (optional) - [Get Docker](https://docs.docker.com/get-docker/)

### 🏃‍♂️ Local Development

1. **Clone the repository**
   ```bash
   git clone https://github.com/vishal210893/Go_Microservice.git
   cd Go_Microservice
   ```

2. **Install dependencies**
   ```bash
   go mod download
   go mod tidy
   ```

3. **Set environment variables**
   ```bash
   export ADDR=:8000
   export DB_ADDR="postgres://username:password@localhost:5432/your_db?sslmode=require"
   export REDIS_ADDR="localhost:6379"
   export REDIS_PW=""
   export SENDGRID_API_KEY="your_sendgrid_api_key"
   export JWT_SECRET="your-super-secret-jwt-key"
   export ENV="development"
   ```

4. **Run database migrations**
   ```bash
   # Install golang-migrate if not already installed
   go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
   
   # Run migrations
   migrate -path cmd/migrate/migrations -database "${DB_ADDR}" up
   ```

5. **Start the application**
   ```bash
   go run cmd/api/main.go
   ```

6. **Verify the installation**
   ```bash
   # Health check
   curl http://localhost:8000/v1/health
   
   # API documentation
   open http://localhost:8000/swagger/index.html
   ```

---

## 🐳 Docker Deployment

### Quick Docker Setup

1. **Pull the pre-built image**
   ```bash
   docker pull vishal210893/go-microservice:1
   ```

2. **Run the container**
   ```bash
   docker run -d \
     --name go-microservice-container \
     -p 8080:8000 \
     -e ADDR=:8000 \
     -e SENDGRID_API_KEY={key} \
     -e REDIS_ADDR=redis-10702.c264.ap-south-1-1.ec2.redns.redis-cloud.com:10702 \
     -e REDIS_PW={pwd} \
     -e DB_ADDR="your_postgresql_connection_string" \
     vishal210893/go-microservice:1
   ```

### Building from Source

1. **Build the Docker image**
   ```bash
   docker build -t go-microservice:latest .
   ```

2. **Run with custom configuration**
   ```bash
   docker run -d \
     --name go-microservice \
     -p 8080:8000 \
     --env-file .env \
     go-microservice:latest
   ```

### Docker Compose (Development)

```bash
# Start all services (app, postgres, redis)
docker-compose up -d

# View logs
docker-compose logs -f app

# Stop all services
docker-compose down
```

### Access Points

| Service | URL | Description |
|---------|-----|-------------|
| **API** | http://localhost:8080 | Main application endpoints |
| **Health Check** | http://localhost:8080/v1/health | Service health status |
| **API Documentation** | http://localhost:8080/swagger/index.html | Interactive API docs |
| **Metrics** | http://localhost:8080/debug/vars | Application metrics |

### Docker Management Commands

```bash
# Container management
docker ps                                    # List running containers
docker logs go-microservice-container       # View application logs
docker exec -it go-microservice-container sh # Access container shell

# Cleanup
docker stop go-microservice-container       # Stop the container
docker rm go-microservice-container         # Remove the container
docker rmi vishal210893/go-microservice:1   # Remove the image
```

---

## ☸️ Kubernetes Deployment

### Prerequisites

- Kubernetes cluster (1.19+)
- `kubectl` configured to access your cluster
- Basic understanding of Kubernetes concepts

### Deployment Steps

1. **Apply ConfigMap for environment variables**
   ```bash
   kubectl apply -f k8s/configmap.yaml
   ```

2. **Create Secrets for sensitive data**
   ```bash
   kubectl apply -f k8s/secret.yaml
   ```

3. **Deploy the application**
   ```bash
   kubectl apply -f k8s/deployment.yaml
   ```

4. **Create service and expose the application**
   ```bash
   kubectl apply -f k8s/service.yaml
   ```

### Kubernetes Configuration Files

```
k8s/
├── configmap.yaml      # Environment configuration
├── secret.yaml         # Sensitive data (API keys, passwords)
├── deployment.yaml     # Application deployment with health checks
└── service.yaml        # LoadBalancer service configuration
```

### Kubernetes Management

```bash
# Check deployment status
kubectl get pods -l app=go-microservice
kubectl get deployments
kubectl get services

# View application logs
kubectl logs -l app=go-microservice -f

# Port forward for local access (development)
kubectl port-forward service/go-microservice-service 8080:80

# Scale the deployment
kubectl scale deployment go-microservice --replicas=3

# Rolling update
kubectl set image deployment/go-microservice go-microservice=vishal210893/go-microservice:2

# Delete resources
kubectl delete -f k8s/

# Port forward for ingress
kubectl port-forward svc/ingress-nginx-controller 8080:80 -n ingress-ngin
```

### Health Checks

The Kubernetes deployment includes:
- **Readiness Probe**: `/v1/health` endpoint
- **Liveness Probe**: `/v1/health` endpoint
- **Resource Limits**: CPU and memory constraints
- **Graceful Shutdown**: SIGTERM handling with 30s timeout

---

## 📁 Project Structure

```
Go_Microservice/
├── cmd/
│   ├── api/                    # API server implementation
│   │   ├── main.go            # Application entry point
│   │   ├── api.go             # Router and middleware setup
│   │   ├── auth.go            # Authentication handlers
│   │   ├── posts.go           # Post management endpoints
│   │   ├── users.go           # User management endpoints
│   │   ├── comments.go        # Comment system endpoints
│   │   ├── health.go          # Health check endpoint
│   │   ├── middleware.go      # Custom middleware
│   │   └── json.go            # JSON utilities
│   └── migrate/               # Database migration tools
│       └── migrations/        # SQL migration files
├── internal/
│   ├── auth/                  # Authentication logic
│   │   ├── jwt.go            # JWT token management
│   │   └── authenticator.go  # Auth interface
│   ├── db/                    # Database connection
│   │   ├── db.go             # PostgreSQL setup
│   │   └── seed.go           # Database seeding
│   ├── repo/                  # Repository layer
│   │   ├── cache/            # Redis caching
│   │   ├── posts.go          # Post repository
│   │   ├── users.go          # User repository
│   │   ├── comments.go       # Comment repository
│   │   └── repository.go     # Repository interfaces
│   ├── mailer/               # Email service
│   │   ├── sendgrid.go       # SendGrid integration
│   │   └── templates/        # Email templates
│   └── ratelimiter/          # Rate limiting
│       └── redis.go          # Redis-based rate limiter
├── docs/                      # Auto-generated API documentation
├── k8s/                       # Kubernetes manifests
├── scripts/                   # Build and deployment scripts
├── docker-compose.yml         # Local development setup
├── Dockerfile                 # Production Docker image
└── Makefile                  # Build automation
```

---

## ⚙️ Configuration

### Environment Variables

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `ADDR` | Server listen address | `:8080` | No |
| `DB_ADDR` | PostgreSQL connection string | - | **Yes** |
| `REDIS_ADDR` | Redis server address | `localhost:6379` | **Yes** |
| `REDIS_PW` | Redis password | - | **Yes** |
| `SENDGRID_API_KEY` | SendGrid API key for emails | - | **Yes** |
| `JWT_SECRET` | JWT signing secret | `secret` | **Yes** |
| `JWT_EXP` | JWT token expiration | `24h` | No |
| `RATE_LIMITER_ENABLED` | Enable rate limiting | `true` | No |
| `RATE_LIMITER_REQUESTS_PER_TIME_FRAME` | Requests per time window | `100` | No |
| `RATE_LIMITER_TIME_FRAME` | Rate limiting time window | `1h` | No |
| `ENV` | Environment (development/production) | `development` | No |

### Example Configuration

```bash
# .env file for local development
ADDR=:8000
DB_ADDR=postgres://user:pass@localhost:5432/microservice_db?sslmode=require
REDIS_ADDR=localhost:6379
REDIS_PW=your_redis_password
SENDGRID_API_KEY=SG.your_sendgrid_key
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production
JWT_EXP=72h
ENV=development
RATE_LIMITER_ENABLED=true
RATE_LIMITER_REQUESTS_PER_TIME_FRAME=100
RATE_LIMITER_TIME_FRAME=1h
```

---

## 📚 API Documentation

### Interactive Documentation

The API documentation is automatically generated using Swagger/OpenAPI and available at:

- **Local Development**: http://localhost:8000/swagger/index.html
- **Docker**: http://localhost:8080/swagger/index.html
- **Production**: https://your-domain.com/swagger/index.html

### Key Endpoints

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| `GET` | `/v1/health` | Service health check | Basic Auth |
| `POST` | `/v1/authentication/user` | Register new user | No |
| `POST` | `/v1/authentication/token` | User login | No |
| `PUT` | `/v1/users/activate/{token}` | Activate user account | No |
| `GET` | `/v1/users/{id}` | Get user profile | JWT |
| `PUT` | `/v1/users/{id}/follow` | Follow user | JWT |
| `PUT` | `/v1/users/{id}/unfollow` | Unfollow user | JWT |
| `GET` | `/v1/users/feed` | Get personalized feed | JWT |
| `POST` | `/v1/posts` | Create new post | JWT |
| `GET` | `/v1/posts/{id}` | Get post details | JWT |
| `PATCH` | `/v1/posts/{id}` | Update post | JWT (Owner/Mod) |
| `DELETE` | `/v1/posts/{id}` | Delete post | JWT (Owner/Admin) |
| `POST` | `/v1/posts/{id}/comments` | Add comment | JWT |

### Authentication

The API uses JWT (JSON Web Tokens) for authentication. Include the token in the Authorization header:

```bash
Authorization: Bearer <your_jwt_token>
```

### Sample API Requests

#### Register a new user
```bash
curl -X POST http://localhost:8000/v1/authentication/user \
  -H "Content-Type: application/json" \
  -d '{
    "username": "johndoe",
    "email": "john@example.com",
    "password": "securepassword123",
    "role_id": 1
  }'
```

#### Create a new post
```bash
curl -X POST http://localhost:8000/v1/posts \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your_jwt_token>" \
  -d '{
    "title": "My First Post",
    "content": "This is the content of my first post!",
    "tags": ["go", "microservice", "api"]
  }'
```
---

## 🔧 Development

### Development Setup

1. **Install development dependencies**
   ```bash
   # Install air for hot reloading
   go install github.com/cosmtrek/air@latest
   
   # Install mockgen for generating mocks
   go install github.com/golang/mock/mockgen@latest
   
   # Install golangci-lint for code quality
   go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
   ```

2. **Start development server with hot reload**
   ```bash
   air
   ```

3. **Database seeding for development**
   ```bash
   go run cmd/api/main.go -seed
   ```

### Code Quality

```bash
# Format code
go fmt ./...

# Run linter
golangci-lint run

# Vet code for common mistakes
go vet ./...

# Generate mocks
go generate ./...
```

### Database Management

```bash
# Create new migration
migrate create -ext sql -dir cmd/migrate/migrations -seq add_new_feature

# Run migrations
migrate -path cmd/migrate/migrations -database "${DB_ADDR}" up

# Rollback migrations
migrate -path cmd/migrate/migrations -database "${DB_ADDR}" down 1

# Check migration status
migrate -path cmd/migrate/migrations -database "${DB_ADDR}" version
```

---

---

## 🐛 Troubleshooting

### Common Issues

#### Database Connection Issues
```bash
# Check database connectivity
pg_isready -h localhost -p 5432 -U your_username

# Verify connection string format
export DB_ADDR="postgres://username:password@host:port/database?sslmode=require"
```

#### Redis Connection Issues
```bash
# Test Redis connection
redis-cli -h your_redis_host -p 6379 ping

# Check Redis authentication
redis-cli -h your_redis_host -p 6379 -a your_password ping
```

#### Docker Issues
```bash
# Check container logs
docker logs go-microservice-container

# Inspect container configuration
docker inspect go-microservice-container

# Access container shell for debugging
docker exec -it go-microservice-container sh
```

#### Kubernetes Issues
```bash
# Check pod status
kubectl describe pod <pod_name>

# View pod logs
kubectl logs <pod_name> -c go-microservice

# Check service endpoints
kubectl get endpoints go-microservice-service
```
---

<div align="center">

**Built with ❤️ using Go**

⭐ **Star this repo if you find it helpful!** ⭐

[Report Bug](https://github.com/vishal210893/Go_Microservice/issues) • [Request Feature](https://github.com/vishal210893/Go_Microservice/issues) • [Contribute](https://github.com/vishal210893/Go_Microservice/pulls)

</div>