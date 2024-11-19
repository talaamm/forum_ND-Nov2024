package main

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)
//Noor
//TOno12345#
func initDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		return nil, err
	}

	err = createTables(db)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func createTables(db *sql.DB) error {
	query := `
    CREATE TABLE IF NOT EXISTS users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        username TEXT NOT NULL UNIQUE,
        email TEXT NOT NULL UNIQUE,
        password TEXT NOT NULL
    );

    CREATE TABLE IF NOT EXISTS posts (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        user_id INTEGER,
        user_name TEXT NOT NULL,
        content TEXT NOT NULL,
        category TEXT NOT NULL,
        postTime TEXT NOT NULL,
        title TEXT NOT NULL,
        image_path TEXT, 
        FOREIGN KEY (user_id) REFERENCES users(id)
    );

    CREATE TABLE IF NOT EXISTS comments (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        post_id INTEGER,
        user_id INTEGER,
        username TEXT,
        content TEXT NOT NULL,
        FOREIGN KEY (post_id) REFERENCES posts(id),
        FOREIGN KEY (user_id) REFERENCES users(id)
    );

	CREATE TABLE IF NOT EXISTS likes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    post_id INTEGER,
    user_id INTEGER,
    is_like BOOLEAN, -- true for like, false for dislike
    FOREIGN KEY (post_id) REFERENCES posts(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);

	CREATE TABLE  IF NOT EXISTS sessions (
    session_id TEXT PRIMARY KEY,    -- A unique identifier for each session (UUID)
    user_id INTEGER,                -- Reference to the user in the users table
    expires_at DATETIME,            -- When the session will expire
    FOREIGN KEY (user_id) REFERENCES users(id)
);

    `
	_, err := db.Exec(query)
	return err
}
