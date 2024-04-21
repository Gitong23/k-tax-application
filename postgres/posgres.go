package postgres

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

type Postgres struct {
	Db *sql.DB
}

func New() (*Postgres, error) {

	// dbUrl := fmt.Sprintf("%s", os.Getenv("DATABASE_URL"))

	// host := os.Getenv("POSTGRES_DB_HOST")
	// port := os.Getenv("POSTGRES_DB_PORT")
	// user := os.Getenv("POSTGRES_DB_USER")
	// password := os.Getenv("POSTGRES_DB_PASSWORD")
	// dbname := os.Getenv("POSTGRES_DB_NAME")

	// // Create database source string
	// databaseSource := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	databaseSource := fmt.Sprintf("%s", os.Getenv("DATABASE_URL"))
	fmt.Println(databaseSource)

	db, err := sql.Open("postgres", databaseSource)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	return &Postgres{Db: db}, nil
}
