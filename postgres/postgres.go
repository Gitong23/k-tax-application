package postgres

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/Gitong23/assessment-tax/config"
	_ "github.com/lib/pq"
)

type Postgres struct {
	Db *sql.DB
}

func New() (*Postgres, error) {
	databaseSource := fmt.Sprintf("%s", config.New().DB.Url)
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
