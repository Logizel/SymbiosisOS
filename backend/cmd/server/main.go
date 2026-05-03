package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	_ "github.com/jackc/pgx/v5/stdlib" // The blank identifier '_' registers the pgx driver with Go's standard sql package

	"symbiosisos/backend/internal/database" // The package sqlc just generated for you
)

func main() {
	// 1. Establish the Database Connection
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		// Fallback for local development. Adjust user/password if yours are different!
		dbURL = "postgres://postgres:postgres@localhost:5432/symbiosisos?sslmode=disable"
	}

	db, err := sql.Open("pgx", dbURL)
	if err != nil {
		log.Fatalf("Fatal error: cannot connect to the database: %v", err)
	}
	defer db.Close() // Ensure the connection pool closes when the app shuts down

	// Ping the database to ensure the connection is actually valid
	if err := db.Ping(); err != nil {
		log.Fatalf("Fatal error: database is unreachable: %v", err)
	}

	// 2. Initialize sqlc queries
	// This creates an instance of the *database.Queries struct, passing in our active connection pool.
	_ = database.New(db) // We assign to '_' for now until we build handlers, but the connection works!

	// 3. Initialize the Chi Router
	router := chi.NewRouter()

	// 4. Attach Global Middleware
	router.Use(middleware.Logger) // Logs every incoming request to the terminal
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"}, // The default port for Vite/React
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
	}))

	// 5. Define a basic Health Check Route
	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("SymbiosisOS Engine is online. Database connected successfully."))
	})

	// 6. Start the HTTP Server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Booting SymbiosisOS server on port %s...", port)
	err = http.ListenAndServe(":"+port, router)
	if err != nil {
		log.Fatal("Server crashed: ", err)
	}
}
