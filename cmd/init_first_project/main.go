package main

import (
	"log"
	"os"

	"github.com/jmoiron/sqlx"
)

func main() {
	if len(os.Args) < 2 {
		log.Printf("Использование: %s <dsn>", os.Args[0])
		log.Printf("Пример: %s \"postgres://postgres:123456@127.0.0.1:5432/dummy\"", os.Args[0])
		os.Exit(1)
	}

	dsn := os.Args[1]

	db, err := Connect(dsn)
	if err != nil {
		log.Fatalf("Connect: %v", err)
	}
	defer func() { _ = db.Close() }()

	if _, err = db.Exec(`
		INSERT INTO projects (name)
		VALUES ('Первая запись')`); err != nil {
		log.Fatalf("вставка первой записи в projects: db.Exec: %v", err)
	}
}

func Connect(dsn string) (*sqlx.DB, error) {
	conn, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err = conn.Ping(); err != nil {
		return nil, err
	}

	return conn, err
}
