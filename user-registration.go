package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func RegisterNewUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		fmt.Println("Invalid request method")
		return
	}

	r.ParseForm()
	pass := r.FormValue("password")
	mail := r.FormValue("email")
	if mail == "" {
		signInErr.Error = "Email is required\n"
		http.Redirect(w, r, "/sign-in", http.StatusSeeOther)
		return
	}
	username := r.FormValue("username")
	if username == "" {
		signInErr.Error = "Missing username\n"
		http.Redirect(w, r, "/sign-in", http.StatusSeeOther)
		return
	}

	if !checkEmail(mail) {
		signInErr.Error = "Invalid email\n"
		http.Redirect(w, r, "/sign-in", http.StatusSeeOther)
		return
	}
	e, ee := checkpassword(pass)
	if !e {
		signInErr.Error = ee + "\n"
		http.Redirect(w, r, "/sign-in", http.StatusSeeOther)
		return
	}
	//check if they already exist
	var userID int
	err := db.QueryRow("SELECT id FROM users WHERE email = ? OR username = ?", mail, username).Scan(&userID)
	if err == nil { // Email or username already exists
		fmt.Println("Email or username already taken")
		signInErr.Error = "Email or username already taken\n"
		http.Redirect(w, r, "/sign-in", http.StatusSeeOther)
		return
	} else if err != sql.ErrNoRows { //ErrNoRows is returned by Scan when QueryRow doesn't return a row.
		fmt.Println("Database error:", err) // Log the error
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to encrypt password", http.StatusInternalServerError)
		return
	}

	//insert user to db
	query := `INSERT INTO users (email, username, password) VALUES (?, ?, ?)`
	_, err = db.Exec(query, mail, username, string(hashedPassword))
	if err != nil {
		log.Printf("Failed to create user: %v", err) // Log the error for debugging
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	signInErr = FormError{Error: "", Success: "Your Account Has been created!\nPlease Log-in Now!"}
	http.Redirect(w, r, "/sign-in", http.StatusSeeOther)
}

func UserLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		fmt.Println("Invalid request method")
		return
	}
	r.ParseForm()
	pass := r.FormValue("password")
	username := r.FormValue("username")

	var storedHashPass string
	err := db.QueryRow("SELECT password FROM users WHERE username = ?", username).Scan(&storedHashPass)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("Database error:", err) // Log the error
			signInErr.Error = "User Doesn't Exist\n"
			http.Redirect(w, r, "/sign-in", http.StatusSeeOther)
			return
		} else {
			http.Error(w, "Server error", http.StatusInternalServerError)
		}
		return
	}

	time.Sleep(500 * time.Millisecond) // Adds a slight delay to prevent timing attacks
	err = bcrypt.CompareHashAndPassword([]byte(storedHashPass), []byte(pass))
	if err != nil {
		// Password doesn't match
		signInErr.Error = "Password incorrect\n"
		http.Redirect(w, r, "/sign-in", http.StatusSeeOther)
		return
	}
	// Set the session token as a cookie
	sessionID := uuid.New().String() // Generate a unique session ID
	cookie := &http.Cookie{
		Name:     "session_token",
		Value:    sessionID,
		Expires:  time.Now().Add(24 * time.Hour), // Set expiration time
		HttpOnly: true,                           // Prevents JavaScript access
		Secure:   true,                           // Requires HTTPS
	}
	http.SetCookie(w, cookie) // Set the cookie in the user's browser
	var userID int
	err = db.QueryRow("SELECT id FROM users WHERE username = ?", username).Scan(&userID)
	if err != nil {
		http.Error(w, "error query sql selecting id by username ", http.StatusUnauthorized)
		return
	}
	_, err = db.Exec(`
	INSERT INTO sessions 
	(session_id, user_id, expires_at) 
	VALUES (?, ?, ?)`, sessionID, userID, cookie.Expires)
	if err != nil {
		http.Error(w, "inserting session error", http.StatusUnauthorized)
		return
	}
	HomePage.User = username
	signInErr.Error = ""
	signInErr.Success = ""
	http.Redirect(w, r, "/home", http.StatusSeeOther)
}

func commentSubmit(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	comment := r.FormValue("comment")
	postIdStr := r.URL.Path[len("/submit-comment/"):] // Gets the part of the URL after "/posted/"
	postID, err := strconv.Atoi(postIdStr)            // Convert the ID to an integer
	if err != nil {
		http.Error(w, "Invalid Post ID", http.StatusBadRequest)
		return
	}
	// Insert the comment into the database
	_, err = db.Exec(`
			INSERT INTO comments (post_id, user_id, content)
			VALUES (?, (SELECT id FROM users WHERE username = ?), ?)
		`, postID, HomePage.User, comment)
	if err != nil {
		http.Error(w, "Error submitting comment: "+err.Error(), http.StatusInternalServerError)
		return
	}
	AllHomePosts, err = fetchHomePosts()
	if err != nil {
		fmt.Println("error fetching home page's posts ", err)
		return
	}
	HomePage.AllPosts = AllHomePosts
	HomePage.AllPosts[len(HomePage.AllPosts)-postID].CommentCount, _ = getCommentCount(postID)
	http.Redirect(w, r, "/posted/"+postIdStr, http.StatusSeeOther)
}
