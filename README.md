# 🚀 Shifty Backend

> **Next-Generation Workforce Management System with AI-Powered Scheduling & Internal Social Network.**

![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat&logo=go)
![Framework](https://img.shields.io/badge/Fiber-v2-black?style=flat&logo=go)
![GraphQL](https://img.shields.io/badge/GraphQL-gqlgen-e535ab?style=flat&logo=graphql)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-Relational-336791?style=flat&logo=postgresql)
![MongoDB](https://img.shields.io/badge/MongoDB-NoSQL-47A248?style=flat&logo=mongodb)
![Redis](https://img.shields.io/badge/Redis-Caching-DC382D?style=flat&logo=redis)
![Qdrant](https://img.shields.io/badge/Qdrant-Vector%20Search-bf0603?style=flat)
![Infrastructure](https://img.shields.io/badge/AWS-Docker-FF9900?style=flat&logo=amazonaws)

## 📖 Introduction

**Shifty** is a comprehensive backend solution designed for the F&B and Retail industries. It moves beyond traditional HRM systems by integrating an **Internal Social Network**, **Real-time Chat**, and **AI Agents** to automate scheduling and improve employee engagement.

The system is built on a **Microservices-ready Monolith** architecture using **Golang (Fiber)**, leveraging a hybrid database approach (Polyglot Persistence) and hybrid API protocols to optimize for speed, scalability, and data integrity.

---

## 🏗 System Architecture

### 1. Polyglot Persistence (Database Strategy)
We use the right tool for the right job:

* **PostgreSQL (Source of Truth):** Handles highly structured data requiring ACID compliance (Users, Shifts, Salaries, Social Graph).
* **MongoDB (High Volume):** Stores unstructured data like Chat History and System Logs (`ai_logs`, `notifications`).
* **Redis (Speed Layer):** Manages Session Caching, Real-time Presence (`user:online`), and expensive calculation caches.
* **Qdrant (Vector DB):** Powers the AI RAG pipeline, storing embeddings for knowledge bases and shift patterns.

### 2. Communication Protocols (Hybrid API)
Shifty implements a **Hybrid API Architecture** to balance performance and flexibility:

* **📡 RESTful API (High Performance):**
    * Used for core resources and state-changing operations.
    * **Use Cases:** Authentication (Auth/Register), File Uploads (Multipart), Webhooks.
    * **Framework:** Go Fiber.

* **⚛️ GraphQL (Flexible Data Fetching):**
    * Used for complex, nested data queries to prevent Over-fetching and Under-fetching.
    * **Use Cases:** Social Newsfeed (Post + User + Comments + Reactions), Shift Assignment Views.
    * **Library:** gqlgen.

* **🔌 gRPC (Internal Communication - Optional):**
    * Used for low-latency communication between the Core Backend and the Python AI Agent Service.

---

## ✨ Key Features

### 📅 Smart Scheduling
* **Multi-tenant:** Support for multiple Restaurant branches.
* **AI Auto-Schedule:** Generates rosters based on historical patterns, `WageMultiplier`, and `IsHoliday` constraints.
* **Shift Management:** Shift swapping, "Open Shift" grabbing, and availability tracking.

### 🌐 Internal Social Network
* **Interactive Feed:** Employees can post updates, images, and announcements.
* **Advanced Comments:** Recursive/Nested comment system (Reply to reply) using Self-Referencing relationships.
* **Reactions:** Rich reaction system (Like, Love, Haha, Angry, etc.) with `OnDelete: CASCADE` logic.

### 💬 Real-time Communication
* **Hybrid Chat:** Relational metadata (Postgres) combined with NoSQL message storage (Mongo).
* **Notifications:** Real-time alerts via WebSockets.

### 📊 HR & Feedback
* **360° Feedback:** Structured review system where Members and Reviewers are clearly defined.
* **Payroll Ready:** Salary stored in `BigInt` (or `Decimal`) for financial precision.

---

## 🛠 Tech Stack

| Category | Technology |
| :--- | :--- |
| **Language** | Go (Golang) 1.22+ |
| **Web Framework** | Fiber v2 |
| **GraphQL** | gqlgen |
| **Databases** | PostgreSQL, MongoDB, Redis, Qdrant |
| **ORM / Drivers** | GORM (SQL), Official Mongo Driver, Go-Redis |
| **Infrastructure** | Docker, Docker Compose, Nginx (Reverse Proxy) |
| **Cloud & CI/CD** | AWS (EC2/S3), GitHub Actions |
| **AI Integration** | OpenAI API / LangChain (via Internal Service) |

---

## 🗄 Database Schema Highlights

The database design implements several advanced patterns:
* **UUID Primary Keys:** Ensures global uniqueness and security.
* **Cascade Deletion:** Logic implemented at the DB level (e.g., Deleting a Post automatically cleans up all related Comments and Reactions).
* **Optional Relationships:** Uses Pointer types (e.g., `*uuid.UUID`) for nullable Foreign Keys (e.g., A User not yet assigned to a Restaurant).

---

## 📂 Project Structure

The project follows **Clean Architecture / Domain-Driven Design (DDD)**:

```bash
shifty-backend/
├── cmd/
│   ├── server/          # Main entry point
│   └── migrate/         # Database migration utility
├── configs/             # Environment variables (.env) & Config structs
├── internal/
│   ├── domain/          # Entities & GORM Models (User, Post, Shift...)
│   ├── repository/      # Database Access Layer (Postgres/Mongo/Redis impl)
│   ├── usecase/         # Business Logic & Service Layer
│   └── delivery/
│       ├── http/        # REST Handlers (Fiber)
│       └── graphql/     # GraphQL Resolvers
├── pkg/
│   ├── database/        # DB Connection Factory
│   ├── logger/          # Structured Logging
│   └── utils/           # Helper functions (JWT, Hash)
├── deploy/              # Dockerfiles & Nginx Configs
└── go.mod

🚀 Getting Started
Prerequisites
Go 1.22+

Docker & Docker Compose

Installation
1. Clone the repository
git clone [https://github.com/your-username/shifty-backend.git](https://github.com/your-username/shifty-backend.git)
cd shifty-backend

2. Environment Setup Create a .env file based on the env example


3. Run Infrastructure (Docker)
docker-compose up -d

4. Run Migrations Initialize the PostgreSQL tables:
go run cmd/migrate/main.go

5. Start Server
go run cmd/server/main.go