package database

import (
	"basic/model"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"os"

	_ "github.com/lib/pq"
)

var Db *sql.DB

func InitializeDb() *sql.DB {
	var err error
	slog.Info("starting................")
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))
	// "postgres://ghiri:develop@postgres-go:5432/items?sslmode=disable"
	slog.Info(connStr)
	Db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	slog.Info("success")

	if err = Db.Ping(); err != nil {
		slog.Error("Postgres ping error : (%v)", err)
	}

	// create items if not there
	query := `
    CREATE TABLE IF NOT EXISTS items (
        id SERIAL PRIMARY KEY,
        name TEXT NOT NULL
    );`
	_, err = Db.Exec(query)

	slog.Info("SUccess")
	if err != nil {
		log.Fatal(err)
	}
	return Db
}

func GetItemsDAO() []model.Item {
	var items []model.Item
	rows, err := Db.Query("SELECT * FROM items")
	if err != nil {
		log.Fatal("FATAL ERROR::::::::::")
		log.Fatal(errors.New("DB error"))
		return nil
	}
	defer rows.Close()

	for rows.Next() {
		var item model.Item
		if err := rows.Scan(&item.ID, &item.Name); err != nil {
			log.Fatal("FATAL ERROR::::::::::")
			log.Fatal(errors.New("Parsing error"))
			return nil
		}
		items = append(items, item)
	}

	return items
}

func CreateItemDAO(newItem model.Item) string {
	err := Db.QueryRow("INSERT INTO items (name) VALUES ($1) RETURNING id;", newItem.Name).Scan(&newItem.ID)
	successMsg := "Successfully Inserted"
	if err != nil {
		log.Fatal("FATAL ERROR::::::::::")
		log.Fatal(errors.New("DB error"))
		successMsg = "Error while inserting"
	}

	return successMsg
}

func UpadteItemDAO(updatedItem model.Item, id string) int64 {
	var rows int64
	result, err := Db.Exec("UPDATE items SET name = $1 WHERE id = $2", updatedItem.Name, updatedItem.ID)
	if err != nil {
		log.Fatal("FATAL ERROR::::::::::")
		log.Fatal(errors.New("DB error"))
		rows = 0
	}
	rows, err = result.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}

	return rows
}

func DeleteItemDAO(id string) int64 {
	var rows int64
	slog.Info("Stating to delete item")
	result, err := Db.Exec("DELETE FROM items WHERE id = $1", id)
	if err != nil {
		log.Fatal("FATAL ERROR::::::::::")
		log.Fatal(errors.New("DB error"))
		rows = 0
	}

	rows, err = result.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}

	return rows
}
