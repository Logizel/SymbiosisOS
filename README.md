# SymbiosisOS

## Overview
SymbiosisOS is a deterministic, logistics-driven B2B marketplace designed to operationalize the circular economy. The platform makes it financially viable for industrial factories to exchange byproducts rather than pay for landfill dumping. 

Unlike fuzzy AI matching systems, SymbiosisOS relies on exact SQL chemical filtering, strict PostGIS geospatial mathematics, and real-time freight calculations to guarantee viable matches based on chemistry, logistics, and compliance.

## The Viability Gate
A match is only successful if it passes three absolute laws of industrial waste:
1. **Chemistry is absolute:** Exact matches on chemical type and purity thresholds.
2. **Freight kills margins:** Geospatial queries ensure the transit distance is viable.
3. **The Arbitrage:** Freight Cost + Processing Cost must be strictly less than the traditional Landfill Cost.

## Architecture & Tech Stack
The system is built as a decoupled monolith optimized for extreme read speeds, complex spatial queries, and absolute type safety.

**Database Engine**
* PostgreSQL
* PostGIS (for native geographic distance calculations via `ST_DWithin`)

**Backend (REST API)**
* Go (Golang)
* Router: `go-chi/chi`
* Database Driver: `jackc/pgx` (native PostGIS support)
* Query Generator: `sqlc` (No ORM, raw SQL compiled to type-safe Go structs)
* Auth: Stateless JWT

**Frontend (Enterprise Portal)**
* React.js (Vite)
* Package Manager: Bun
* Styling: Tailwind CSS + Shadcn UI
* State Management: Zustand
* Visualization: Recharts

---

## Prerequisites
Before you begin, ensure you have the following installed:
* [Go](https://golang.org/doc/install) (1.21+)
* [Bun](https://bun.sh/)
* [PostgreSQL](https://www.postgresql.org/download/) with the [PostGIS extension](https://postgis.net/install/)
* [sqlc](https://docs.sqlc.dev/en/latest/tutorials/getting-started.html) (`go install [github.com/sqlc-dev/sqlc/cmd/sqlc@latest](https://github.com/sqlc-dev/sqlc/cmd/sqlc@latest)`)

---

## Local Development Setup

### 1. Database Setup
Ensure PostgreSQL is running and create the database with the PostGIS extension.

```sql
CREATE DATABASE symbiosisos;
\c symbiosisos
CREATE EXTENSION postgis;
```
Execute your `schema.sql` files against this database to create the `users`, `facilities`, `waste_streams`, and `buyer_requirements` tables.

### 2. Backend Setup
Navigate to the project root and start the Go backend.

```bash
# Initialize the backend directory
mkdir backend && cd backend
go mod init symbiosisos/backend

# Install core Go dependencies
go get github.com/go-chi/chi/v5
go get github.com/jackc/pgx/v5
go get github.com/golang-jwt/jwt/v5
go get github.com/go-chi/cors

# Generate type-safe database models from your SQL queries
sqlc generate

# Start the Go server
go run cmd/server/main.go
```

### 3. Frontend Setup
Open a new terminal window from the project root to set up the React client.

```bash
# Initialize the frontend using Bun and Vite
bun create vite frontend --template react-ts
cd frontend

# Install core frontend dependencies
bun add zustand recharts tailwindcss postcss autoprefixer
bun add -d @types/node

# Start the development server
bun run dev
```

---

## Project Structure

```text
symbiosisos/
├── backend/
│   ├── cmd/
│   │   └── server/
│   │       └── main.go          # Application entry point
│   ├── internal/
│   │   ├── auth/                # JWT validation and generation
│   │   ├── database/            # sqlc generated Go structs (Do not edit directly)
│   │   └── handlers/            # HTTP route handlers (Business logic)
│   ├── sql/
│   │   ├── queries/             # Raw SQL SELECT/INSERT/UPDATE files
│   │   └── schema/              # Raw SQL table definitions
│   ├── sqlc.yaml                # sqlc configuration file
│   └── go.mod                   # Go module definitions
└── frontend/
    ├── src/
    │   ├── components/          # Shadcn UI and generic components
    │   ├── features/            # Domain-specific React components (Waste Forms, Matches)
    │   ├── store/               # Zustand state stores
    │   └── App.tsx              # Root React component
    ├── package.json             # Bun dependencies
    └── vite.config.ts           # Vite configuration
```
