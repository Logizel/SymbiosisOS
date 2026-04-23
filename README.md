# SymbiosisOS

An operating system for the circular economy. SymbiosisOS is a B2B platform designed to solve the industrial waste crisis by connecting "Waste Generators" (factories producing waste) with "Material Buyers" (facilities that can use that waste as raw material).

Instead of paying high fees to dump materials in landfills, SymbiosisOS uses an AI-powered matching engine to find profitable and compliant exchange opportunities.

**Repository:** https://github.com/Logizel/SymbiosisOS

## How It Works

1. **Upload:** A factory uploads a "Waste Passport" (a PDF detailing their waste stream).
2. **Extract:** The backend parses the document using AI to extract key data like chemical composition, quantity, and disposal costs.
3. **Match:** The system compares the waste stream against buyer requirements in the database using vector similarity and logistics calculations.
4. **Connect:** Users view their matches, authorize agreements, and the system generates a cryptographic hash for an immutable tracking certificate.

## Tech Stack

### Frontend

- React.js (Vite)
- Tailwind CSS & Shadcn UI
- Zustand (State management)
- Recharts (Data visualization)

### Backend & Database

- ElysiaJS (API framework)
- PostgreSQL
- Prisma ORM
- pgvector (For vector-based semantic matching)

### Intelligence Layer

- Document Parsing (LlamaParse / pdf-parse)
- LLM (For structuring extracted data)

### Tooling

- Bun (Package manager and runtime)

## Local Setup

Make sure you have Bun and PostgreSQL installed on your machine.

### 1. Clone the repository

```bash
git clone https://github.com/Logizel/SymbiosisOS
cd SymbiosisOS
```

### 2. Install dependencies

```bash
bun install
```

### 3. Environment setup

Create a `.env` file in the root directory and add your database URL and necessary API keys.

```plaintext
DATABASE_URL="postgresql://user:password@localhost:5432/symbiosisos"
OPENROUTER_API_KEY="your_api_key_here"
```

### 4. Setup the database

Push the Prisma schema to your PostgreSQL database.

```bash
bunx prisma db push
```

### 5. Start the development server

```bash
bun run dev
```
