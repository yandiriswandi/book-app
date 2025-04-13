package config

import (
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var (
	DB  *sqlx.DB
	err error
)

func InitDB() {

	err = godotenv.Load(".env")
	if err != nil {
		panic("Error loading .env file")
	}
	var err error
	// Ambil value dari environment
	// dbHost := os.Getenv("DB_HOST")
	// dbPort := os.Getenv("DB_PORT")
	// dbUser := os.Getenv("DB_USER")
	// dbPassword := os.Getenv("DB_PASSWORD")
	// dbName := os.Getenv("DB_NAME")

	dbHost := os.Getenv("PGHOST")
	dbPort := os.Getenv("PGPORT")
	dbUser := os.Getenv("PGUSER")
	dbPassword := os.Getenv("PGPASSWORD")
	dbName := os.Getenv("PGDATABASE")

	// Bangun connection string
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName,
	)

	DB, err = sqlx.Connect("postgres", connStr)
	if err != nil {
		log.Fatal("Gagal koneksi ke database:", err)
	}

	// schema := `
	// CREATE TABLE IF NOT EXISTS bioskop (
	// 	id SERIAL PRIMARY KEY,
	// 	nama VARCHAR(100) NOT NULL,
	// 	lokasi VARCHAR(100) NOT NULL,
	// 	rating FLOAT
	// );`
	schema := `
	CREATE TABLE IF NOT EXISTS category (
		id SERIAL PRIMARY KEY,
		nama VARCHAR(100) NOT NULL,
		description TEXT
	);
	
	CREATE TABLE IF NOT EXISTS bioskop (
		id SERIAL PRIMARY KEY,
		nama VARCHAR(100) NOT NULL,
		kategori_id INT NOT NULL,
		description TEXT,
		rating FLOAT,
		FOREIGN KEY (kategori_id) REFERENCES category(id)
	);`

	DB.MustExec(schema)
	fmt.Println("Database siap digunakan 🚀")
}
