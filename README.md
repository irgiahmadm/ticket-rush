# **TicketRush 🎟️**

TicketRush is a high-performance ticketing system demonstrating the transition from a Monolith to a Microservices architecture using the **Strangler Fig Pattern**. It implements **Hexagonal Architecture** (Ports & Adapters) in Go.

## **🏗️ Architecture**

- **Pattern:** Strangler Fig (API Gateway routes traffic between Microservices and the Legacy Monolith).
- **Architecture:** Hexagonal (Ports & Adapters) for loose coupling.
- **Gateway:** Go Chi with Redis Rate Limiting.
- **Services:**
  - **Auth Service:** Go \+ Gin \+ PostgreSQL (Users, JWT).
  - **Order Service:** Go \+ Gin \+ PostgreSQL (Orders).
  - **Monolith:** Legacy fallback service.

## **🛠️ Tech Stack**

- **Language:** Go (Golang) 1.25+
- **Frameworks:** Gin (Services), Chi (Gateway)
- **Database:** PostgreSQL (Local Instance)
- **Caching/Rate Limiting:** Redis
- **Config:** Viper
- **Containerization:** Docker & Docker Compose

## **🚀 Prerequisites**

1. **Go:** Installed locally for module initialization.
2. **Docker & Docker Compose:** For running the services and Redis.
3. **PostgreSQL:** A local installation running on port 5432\.
   - _Note: This setup relies on a local DB, not a Dockerized one._

## **⚙️ Setup & Configuration**

### **1\. Database Setup**

Ensure your local PostgreSQL is running and accessible.

1. Create a database named ticketrush.
2. Run the SQL commands from init.sql to create tables and seed the admin user.
3. **Critical:** Ensure your pg_hba.conf allows connections from Docker (see Troubleshooting below).

### **2\. Environment Variables**

Each service has its own .env file. Ensure DATABASE_URL points to host.docker.internal.  
**auth-service/.env & order-service/.env:**  
```
PORT=3001 \# (3002 for order-service)  
\# Replace 'postgres:password' with your local DB credentials  
DATABASE_URL=postgres://postgres:yourpassword@host.docker.internal:5432/ticketrush?sslmode=disable  
JWT_SECRET=supersecretkey
```

**gateway/.env:**  
```
PORT=3000  
REDIS_ADDR=redis:6379  
JWT_SECRET=supersecretkey  
AUTH_SERVICE_URL=http://auth-service:3001  
ORDER_SERVICE_URL=http://order-service:3002  
MONOLITH_URL=http://monolith_service:8080  
RATE_LIMIT_REQ=10  
RATE_LIMIT_WINDOW=60
```

### **3\. Initialize Modules**

If you haven't already, initialize Go modules for each service:  
```
cd gateway && go mod tidy && cd ..  
cd auth-service && go mod tidy && cd ..  
cd order-service && go mod tidy && cd ..
```

## **▶️ How to Run**

Start the entire stack using Docker Compose:  
`docker-compose up --build`

This will start:

- **Gateway:** http://localhost:3000
- **Auth Service:** Internal port 3001
- **Order Service:** Internal port 3002
- **Redis:** Port 6379
- **Monolith:** Internal port 8080

## **🧪 API Endpoints**

### **1\. Register (Auth Service)**
```
curl \-X POST http://localhost:3000/auth/register \\  
 \-H "Content-Type: application/json" \\  
 \-d '{"email": "newuser@example.com", "password": "password123"}'
```

### **2\. Login (Auth Service)**

_Default Admin: admin@example.com / password_  
```
curl \-X POST http://localhost:3000/auth/login \\  
 \-H "Content-Type: application/json" \\  
 \-d '{"email": "admin@example.com", "password": "password"}'
```

_Copy the token from the response._

### **3\. Create Order (Order Service)**

_Requires Bearer Token. Replace \<UUID\> with a valid User ID from your DB._  
```
curl \-X POST http://localhost:3000/orders \\  
 \-H "Authorization: Bearer \<YOUR_TOKEN\>" \\  
 \-H "Content-Type: application/json" \\  
 \-d '{"user_id": "\<UUID_FROM_DB\>", "product_id": 101}'
```

### **4\. Legacy Monolith Fallback**

Any route not matched by microservices falls back here.  
curl http://localhost:3000/legacy/hello

## **🐛 Troubleshooting**

**dial tcp \[::1\]:5432: connect: connection refused**

- **Cause:** The container is trying to connect to itself via localhost.
- **Fix:** Ensure .env uses host.docker.internal and docker-compose.yml has extra_hosts configured.

**no pg_hba.conf entry for host ...**

- **Cause:** Local Postgres is blocking the Docker connection.
- **Fix:** Edit your local pg_hba.conf and add:  
  host all all 0.0.0.0/0 md5

  Then restart the PostgreSQL service.

**invalid input syntax for type uuid: ""**

- **Cause:** Sending an empty string for user_id.
- **Fix:** Ensure you copy a real UUID from the users table in your database and put it in the request JSON.
