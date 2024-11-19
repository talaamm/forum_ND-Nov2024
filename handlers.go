package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"
)

func mypostsHandler(w http.ResponseWriter, r *http.Request) {
	userid, err := getUserFromSession(w, r, db)
	if err != nil {
		http.Error(w, "err getting userid from session", http.StatusBadRequest)
		return
	}
	HomePage.AllPosts, err = myPostsFilter(userid)
	if err != nil {
		http.Error(w, "err getting myposts", http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, "/home", http.StatusSeeOther)
}

func myPostsFilter(user_id int) ([]postsForHome, error) {
	rows, err := db.Query("SELECT id, user_name, title FROM posts WHERE user_id = ?", user_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []postsForHome
	for rows.Next() {
		var post postsForHome
		err := rows.Scan(&post.PostId, &post.Username, &post.Title)
		if err != nil {
			return nil, err
		}
		post.LikeCount, err = getLikeCount(post.PostId)
		if err != nil {
			return nil, err
		}
		post.DisLikeCount, err = getDislikeCount(post.PostId)
		if err != nil {
			return nil, err
		}
		post.CommentCount, err = getCommentCount(post.PostId)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func getPostsByCategory(category string) ([]postsForHome, error) {
	var rows *sql.Rows
	var err error
	// If the category is "All", we want to return all posts
	if category == "All" {
		rows, err = db.Query(`	SELECT id, user_name, title
		FROM posts 
		ORDER BY id DESC;`)

	} else {
		rows, err = db.Query("SELECT id, user_name, title FROM posts WHERE category = ?", category)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []postsForHome
	for rows.Next() {
		var post postsForHome
		err := rows.Scan(&post.PostId, &post.Username, &post.Title)
		if err != nil {
			return nil, err
		}
		post.LikeCount, err = getLikeCount(post.PostId)
		if err != nil {
			return nil, err
		}
		post.DisLikeCount, err = getDislikeCount(post.PostId)
		if err != nil {
			return nil, err
		}
		post.CommentCount, err = getCommentCount(post.PostId)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func filterPostsHandler(w http.ResponseWriter, r *http.Request) {
	category := r.URL.Path[len("/category/"):] // Gets the part of the URL after "/posted/"
	if category == "" {
		http.Error(w, "Category not specified", http.StatusBadRequest)
		return
	}
	fmt.Println(category)

	// Fetch posts by the selected category
	posts, err := getPostsByCategory(category)
	if err != nil {
		http.Error(w, "Error fetching posts", http.StatusInternalServerError)
		return
	}
	HomePage.AllPosts = posts
	http.Redirect(w, r, "/home", http.StatusSeeOther)

}

func homeGuest(guest_posts []postsForHome) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		GuestPost := data{AllPosts: guest_posts}

		err := templates.ExecuteTemplate(w, "index.html", GuestPost)
		if err != nil {
			log.Printf("Template execution error: %v", err)
			http.Error(w, "Unable to render template", http.StatusInternalServerError)
			return
		}
	}
}

func loginPage(e *FormError) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(e)
		err := templates.ExecuteTemplate(w, "sign-in.html", e)
		if err != nil {
			http.Error(w, "Unable to render template", http.StatusInternalServerError)
		}
	}
}

func homePageTmpl(w http.ResponseWriter, r *http.Request) {
	fmt.Println(HomePage.User)
	fromWhere = ""
	err := templates.ExecuteTemplate(w, "myhome.html", HomePage)
	if err != nil {
		log.Printf("Failed to render template: %v", err) // Log the error for debugging
		http.Error(w, "Unable to render template", http.StatusInternalServerError)
		return
	}
}

func logOutHandler(w http.ResponseWriter, r *http.Request) {

	cookie, err := r.Cookie("session_token") //This retrieves the session cookie from the request to get the session_token.
	if err != nil {
		http.Error(w, "Session cookie not found", http.StatusBadRequest)
		return
	}
	sessionID := cookie.Value
	err = deleteSession(sessionID, db)
	if err != nil {
		http.Error(w, "Failed to delete session", http.StatusInternalServerError)
		return
	}
	cookie = &http.Cookie{
		Name:     "session_token",
		Value:    "",                             // Clear the value
		Expires:  time.Now().Add(-1 * time.Hour), // Expire immediately
		Path:     "/",                            // Make sure the Path is set to root to delete the cookie
		HttpOnly: true,                           // Optional, for security, ensures cookie is not accessible via JavaScript
	}
	http.SetCookie(w, cookie) // Send the cookie back to the browser

	// Redirect user to login or homepage after logging out
	http.Redirect(w, r, "/sign-in", http.StatusSeeOther)
}
