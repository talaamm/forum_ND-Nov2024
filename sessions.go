package main

import (
	"database/sql"
	"errors"
	"net/http"
	"time"
)

// GetUserFromSession retrieves the user associated with the given session token (cookie).
func getUserFromSession(w http.ResponseWriter, r *http.Request, db *sql.DB) (int, error) {
	// Retrieve the session cookie
	cookie, err := r.Cookie("session_token")
	if err != nil {
		// No session token in the request
		if errors.Is(err, http.ErrNoCookie) {
			http.Error(w, "Session not found", http.StatusUnauthorized)
			return 0, errors.New("no session cookie")
		}
		// Other cookie error
		http.Error(w, "Failed to retrieve session cookie", http.StatusBadRequest)
		return 0, err
	}

	sessionID := cookie.Value

	// Fetch session from the database
	var userID int
	var expiresAt time.Time
	err = db.QueryRow("SELECT user_id, expires_at FROM sessions WHERE session_id = ?", sessionID).Scan(&userID, &expiresAt)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Invalid session", http.StatusUnauthorized)
			return 0, errors.New("session not found")
		}
		http.Error(w, "Server error", http.StatusInternalServerError)
		return 0, err
	}

	// Check if the session has expired
	if time.Now().After(expiresAt) {
		// Session has expired, delete it from the database
		_, err := db.Exec("DELETE FROM sessions WHERE session_id = ?", sessionID)
		if err != nil {
			http.Error(w, "Failed to remove expired session", http.StatusInternalServerError)
			return 0, err
		}
		http.Error(w, "Session has expired", http.StatusUnauthorized)
		return 0, errors.New("session expired")
	}

	// Session is valid, return the user ID
	return userID, nil
}

func deleteSession(sessionID string, db *sql.DB) error {
	_, err := db.Exec("DELETE FROM sessions WHERE session_id = ?", sessionID)
	return err
}
