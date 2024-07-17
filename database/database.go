package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	db *sql.DB
}

func NewDatabase(dbPath string) (*Database, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS subscribers (user_id INTEGER PRIMARY KEY)`)
	if err != nil {
		return nil, err
	}

	return &Database{db: db}, nil
}

func (d *Database) Close() error {
	return d.db.Close()
}

func (d *Database) Subscribe(userID int64) error {
	_, err := d.db.Exec("INSERT OR IGNORE INTO subscribers (user_id) VALUES (?)", userID)
	return err
}

func (d *Database) Unsubscribe(userID int64) error {
	_, err := d.db.Exec("DELETE FROM subscribers WHERE user_id = ?", userID)
	return err
}

func (d *Database) GetSubscribers() ([]int64, error) {
	rows, err := d.db.Query("SELECT user_id FROM subscribers")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subscribers []int64
	for rows.Next() {
		var userID int64
		err := rows.Scan(&userID)
		if err != nil {
			log.Printf("Error scanning user ID: %v", err)
			continue
		}
		subscribers = append(subscribers, userID)
	}

	return subscribers, nil
}
