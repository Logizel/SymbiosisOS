package main

import (
	"database/sql"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	_ "github.com/jackc/pgx/v5/stdlib"
	"log"
	"net/http"
	"os"

	"symbiosisos/backend/internal/database"
	"symbiosisos/backend/internal/handlers"
)

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5432/symbiosisos?sslmode=disable"
	}

	db, err := sql.Open("pgx", dbURL)
	if err != nil {
		log.Fatalf("Fatal error: cannot connect to the database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Fatal error: database is unreachable: %v", err)
	}

	dbQueries := database.New(db)
	apiCfg := &handlers.APIConfig{
		DB: dbQueries,
	}

	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
	}))

	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		handlers.RespondWithJSON(w, http.StatusOK, map[string]string{"status": "online"})
	})

	router.Route("/api/v1", func(r chi.Router) {
		r.Post("/users", apiCfg.HandlerCreateUser)
		r.Post("/login", apiCfg.HandlerLogin)

		r.With(apiCfg.MiddlewareAuth).Get("/users/me", apiCfg.HandlerGetMe)
		r.With(apiCfg.MiddlewareAuth).Post("/facilities", apiCfg.HandlerCreateFacility)
		r.With(apiCfg.MiddlewareAuth).Post("/waste", apiCfg.HandlerCreateWasteStream)
		r.With(apiCfg.MiddlewareAuth).Post("/requirements", apiCfg.HandlerCreateRequirement)
		r.With(apiCfg.MiddlewareAuth).Get("/matches/{facility_id}", apiCfg.HandlerGetMatches)
		r.With(apiCfg.MiddlewareAuth).Post("/transactions", apiCfg.HandlerCreateTransaction)
		r.With(apiCfg.MiddlewareAuth).Get("/facilities", apiCfg.HandlerGetFacilities)
		r.With(apiCfg.MiddlewareAuth).Get("/waste/{facility_id}", apiCfg.HandlerGetFacilityWaste)
		r.With(apiCfg.MiddlewareAuth).Get("/requirements/{facility_id}", apiCfg.HandlerGetFacilityRequirements)
		r.With(apiCfg.MiddlewareAuth).Get("/transactions", apiCfg.HandlerGetTransactions)
		r.With(apiCfg.MiddlewareAuth).Delete("/waste/{id}", apiCfg.HandlerDeleteWasteStream)
		r.With(apiCfg.MiddlewareAuth).Delete("/requirements/{id}", apiCfg.HandlerDeleteRequirement)
	})

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
