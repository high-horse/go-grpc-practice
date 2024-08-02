package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	sqlc "grpc-1/store/database"
)

type Config struct {
	Driver   string
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func LoadConfig() Config {
	return Config{
		Driver:   getEnv("DB_DRIVER", "postgres"),
		Host:     getEnv("DB_HOST", "postgres"),
		Port:     getEnvInt("DB_PORT", 5432),
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "root"),
		DBName:   getEnv("DB_NAME", "fiber"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}
}

func getEnv(key, default_key string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return default_key
}

func getEnvInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		intValue, err := strconv.Atoi(value)
		if err == nil {
			return intValue
		}
	}
	return defaultValue
}

var (
	DB      *sql.DB
	Queries *sqlc.Queries
)

func ConnectDB() error {
	config := LoadConfig()

	dsn := fmt.Sprintf(
		"%s://%s:%s@%s:%d/%s?sslmode=%s",
		config.Driver, config.User, config.Password, config.Host, config.Port, config.DBName, config.SSLMode,
	)
	db, err := sql.Open("postgres", dsn)
	// db, err := sql.Open(config.Driver, dsn)
	if err != nil {
		return fmt.Errorf("couldnot connect to DB : %w", err)
	}

	if err := db.Ping(); err != nil {
		return fmt.Errorf("couldnot ping to DB : %w", err)
	}
	DB = db
	Queries = sqlc.New(DB)

	log.Println("Connected database successfully.")
	return nil
}

func DisConnectDB() {
	if DB != nil {
		DB.Close()
		log.Println("Database Connectino Closed")
	}
}
