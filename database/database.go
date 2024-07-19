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

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			user_id INTEGER PRIMARY KEY,
			authorized INTEGER NOT NULL DEFAULT 0,
			admin INTEGER NOT NULL DEFAULT 0,
			subscribed INTEGER NOT NULL DEFAULT 0
		)
	`)
	if err != nil {
		return nil, err
	}

	return &Database{db: db}, nil
}

func (d *Database) Close() error {
	return d.db.Close()
}

func (d *Database) AddNewUser(userId int64) error {
	_, err := d.db.Exec("INSERT OR IGNORE INTO users (user_id) VALUES (?)", userId)
	return err
}

func (d *Database) IsAuthorized(userId int64) (int, error) {
	var match int
	err := d.db.QueryRow("SELECT CASE WHEN authorized = 1 THEN 1 ELSE 0 END FROM users WHERE user_id = ?", userId).Scan(&match)
	return match, err
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
