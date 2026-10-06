package db

import (
	"database/sql"
	"embed"
	"io/fs"
	"log"
	"os"
	"sort"
	"time"

	_ "github.com/lib/pq"
)

// migrations are compiled into the binary, so they work regardless of the working directory
//
//go:embed migrations/*.sql
var migrationsFS embed.FS

func InitDB() *sql.DB {
	db, err := sql.Open("postgres", os.Getenv("DB_URL"))
	if err != nil {
		log.Fatal(err)
	}

	db.SetConnMaxIdleTime(5)
	db.SetConnMaxLifetime(time.Minute * 5)
	db.SetMaxOpenConns(30)

	waitForDB(db)
	runMigrations(db)

	return db
}

// waitForDB retries until postgres accepts connections (sql.Open does not connect)
func waitForDB(db *sql.DB) {
	var err error
	for range 30 {
		if err = db.Ping(); err == nil {
			return
		}
		log.Printf("waiting for database: %v", err)
		time.Sleep(time.Second)
	}
	log.Fatalf("database not reachable: %v", err)
}

// runMigrations executes every .sql file in migrations/ in filename order
func runMigrations(db *sql.DB) {
	files, err := fs.Glob(migrationsFS, "migrations/*.sql")
	if err != nil {
		log.Fatal(err)
	}
	sort.Strings(files)

	for _, file := range files {
		query, err := migrationsFS.ReadFile(file)
		if err != nil {
			log.Fatal(err)
		}
		if _, err := db.Exec(string(query)); err != nil {
			log.Fatalf("migration %s failed: %v", file, err)
		}
		log.Printf("migration %s applied", file)
	}
}
