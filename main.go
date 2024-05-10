package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"p-system/repositories/transaction"
	"p-system/repositories/user"
	"p-system/repositories/wallet"
	"p-system/services/transactionsservice"
	"time"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	DB "p-system/db"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
	userName := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	dbName := os.Getenv("POSTGRES_DB")

	// Set up database connection
	dbUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", userName, password, host, port, dbName)

	fmt.Println(dbUrl)

	db, err := sqlx.Connect("postgres", dbUrl)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer db.Close()

	// Ping database to verify connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("Error pinging database: %v", err)
	}

	// Run Goose migrations
	if err := DB.Migrate(db.DB); err != nil {
		log.Fatalf("Error running migrations: %v", err)
	}
	log.Println("Migrations ran successfully")

	// Initialize service with repositories and other dependencies
	userRepo := user.NewRepository(db)
	walletRepo := wallet.NewRepository(db)
	transactionRepo := transaction.NewRepository(db)
	svc := transactionsservice.NewService(userRepo, transactionRepo, walletRepo)

	// Create a new router
	r := mux.NewRouter()

	// Define routes
	r.HandleFunc("/transactions", svc.HandleTransaction).Methods("POST")

	// Create a server instance
	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	log.Println("Server started on port 8080")
	log.Fatal(server.ListenAndServe())
}
