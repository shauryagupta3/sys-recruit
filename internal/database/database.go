package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"recruit-sys/internal/models"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/joho/godotenv/autoload"
)

// Service represents a service that interacts with a database.
type Service interface {
	CreateUser(*models.User) error
	SelectUserWhereMail(string) (models.User, error)
	SelectUserWhereID(float64) (models.User, error)

	CreateJob(*models.Job) error
	SelectAllJobs() ([]models.Job, error)
	SelectJobsPostedBy(float64) ([]models.Job, error)
	SelectJobByIdAdmin(float64,int) (models.Job, error)
	SelectJobsAppliedBy(float64) ([]models.Job, error)
	SelectJobsByID(int) (models.Job, error)

	CreateProfile(*models.Profile) error
	SelectProfileById(float64) (models.Profile,error)
	SelectAllProfiles() ([]models.Profile,error)

	ApplyToJob(JobId int, UserID int) error
	// Health returns a map of health status information.
	// The keys and values in the map are service-specific.
	Health() map[string]string

	// Close terminates the database connection.
	// It returns an error if the connection cannot be closed.
	Close() error
}

type service struct {
	db *pgxpool.Pool
}

var (
	database   = os.Getenv("DB_DATABASE")
	password   = os.Getenv("DB_PASSWORD")
	username   = os.Getenv("DB_USERNAME")
	port       = os.Getenv("DB_PORT")
	host       = os.Getenv("DB_HOST")
	schema     = os.Getenv("DB_SCHEMA")
	dbInstance *service
)

func testConnection(ctx context.Context, db *pgxpool.Pool) error {
	var now time.Time
	err := db.QueryRow(ctx, "SELECT NOW()").Scan(&now)
	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}
	fmt.Println("Current time from database:", now)
	return nil
}

func New() Service {
	// Reuse Connection
	if dbInstance != nil {
		return dbInstance
	}
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&search_path=%s", username, password, host, port, database, schema)
	_, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	db, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	if err := testConnection(context.Background(), db); err != nil {
		log.Fatal("Database connection test failed:", err)
		os.Exit(1)
	}

	fmt.Println("database connected")

	// if err := DropTables(context.Background(), db); err != nil {
	// 	log.Fatal("Failed to create tables : ", err)
	// }

	if err := CreateTables(context.Background(), db); err != nil {
		log.Fatal("Failed to create tables : ", err)
	}

	dbInstance = &service{
		db: db,
	}

	return dbInstance
}

// Health checks the health of the database connection by pinging the database.
// It returns a map with keys indicating various health statistics.
func (s *service) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	stats := make(map[string]string)

	// Ping the database
	err := s.db.Ping(ctx)
	if err != nil {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("db down: %v", err)
		log.Fatalf(fmt.Sprintf("db down: %v", err)) // Log the error and terminate the program
		return stats
	}

	// Database is up, add more statistics
	stats["status"] = "up"
	stats["message"] = "It's healthy"

	// Get database stats (like open connections, in use, idle, etc.)
	dbStats := s.db.Stat()
	stats["total_connections"] = strconv.Itoa(int(dbStats.TotalConns()))
	stats["idle_connections"] = strconv.Itoa(int(dbStats.IdleConns()))
	stats["used_connections"] = strconv.Itoa(int(dbStats.AcquiredConns()))
	stats["max_connections"] = strconv.Itoa(int(dbStats.MaxConns()))
	stats["acquire_count"] = strconv.FormatInt(dbStats.AcquireCount(), 10)
	stats["acquire_duration"] = dbStats.AcquireDuration().String()
	stats["canceled_acquire_count"] = strconv.FormatInt(dbStats.CanceledAcquireCount(), 10)

	// Evaluate stats to provide a health message
	if dbStats.TotalConns() > 40 { // Assuming 50 is the max for this example
		stats["message"] = "The database is experiencing heavy load."
	}

	if dbStats.AcquireCount() > 1000 {
		stats["message"] = "The database has a high number of wait events, indicating potential bottlenecks."
	}

	if dbStats.CanceledAcquireCount() > int64(dbStats.TotalConns())/2 {
		stats["message"] = "Many acquire attempts are being canceled, consider revising the connection pool settings."
	}

	if dbStats.AcquireDuration() > time.Hour {
		stats["message"] = "Many connections are being held for a long duration, consider optimizing query performance or increasing pool size."
	}

	return stats
}

// Close closes the database connection.
// It logs a message indicating the disconnection from the specific database.
// If the connection is successfully closed, it returns nil.
// If an error occurs while closing the connection, it returns the error.
func (s *service) Close() error {
	log.Printf("Disconnected from database: %s", database)
	s.db.Close()
	return nil
}
