package main

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func likePost(w http.ResponseWriter, r *http.Request) {
	userid, err := getUserFromSession(w, r, db)
	if err != nil {
		http.Error(w, "Error getUserFromSession", http.StatusInternalServerError)
		return
	}
	postID := r.FormValue("post_id")
	isLike := r.FormValue("like") == "true"
	var existingLike bool
	err = db.QueryRow("SELECT is_like FROM likes WHERE post_id = ? AND user_id = ?", postID, userid).Scan(&existingLike)
	if err != nil && err != sql.ErrNoRows {
		http.Error(w, "Error checking like status", http.StatusInternalServerError)
		return
	}

	if err == sql.ErrNoRows {
		// No existing like/dislike, insert a new record
		_, err = db.Exec("INSERT INTO likes (post_id, user_id, is_like) VALUES (?, ?, ?)", postID, userid, isLike)
		if err != nil {
			http.Error(w, "Error liking/disliking post", http.StatusInternalServerError)
			return
		}
	} else {
		// Update the existing like/dislike record
		_, err = db.Exec("UPDATE likes SET is_like = ? WHERE post_id = ? AND user_id = ?", isLike, postID, userid)
		if err != nil {
			http.Error(w, "Error updating like/dislike", http.StatusInternalServerError)
			return
		}
	}
	idd, _ := strconv.Atoi(postID)
	var maxID int
	if err := db.QueryRow("SELECT MAX(id) FROM posts").Scan(&maxID); err != nil {
		log.Fatal(err)
	}
	AllHomePosts, err = fetchHomePosts()
	if err != nil {
		fmt.Println("error fetching home page's posts ", err)
		return
	}
	HomePage.AllPosts = AllHomePosts
	HomePage.AllPosts[maxID-idd].LikeCount, _ = getLikeCount(idd)
	HomePage.AllPosts[maxID-idd].DisLikeCount, _ = getDislikeCount(idd)
	http.Redirect(w, r, "/posted/"+postID, http.StatusSeeOther) // Redirect back to the homepage or the post
}

func getCommentCount(postID int) (int, error) {
	var commentCount int
	err := db.QueryRow("SELECT COUNT(*) FROM comments WHERE post_id = ?", postID).Scan(&commentCount)
	if err != nil {
		return 0, err
	}
	return commentCount, nil
}

func creatApost(w http.ResponseWriter, r *http.Request) {
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	if fromWhere == "" {
		postCreateError = ""
	}
	data := postCreate{Username: HomePage.User, DateNow: currentTime, Err: postCreateError}
	err := templates.ExecuteTemplate(w, "create.html", data)
	if err != nil {
		http.Error(w, "Unable to render template", http.StatusInternalServerError)
		return
	}
}

func postPostTodb(w http.ResponseWriter, r *http.Request) {
	userid, err := getUserFromSession(w, r, db) // retrieve logged-in user
	if err != nil {
		return // Error is already handled in getUserFromSession
	}
	timee := r.FormValue("postTime")
	cat := r.FormValue("postCategory")
	cont := r.FormValue("postContent")
	ttl := r.FormValue("postTitle")
	var userName string
	err = db.QueryRow("SELECT username FROM users WHERE id = ?", userid).Scan(&userName)
	if err != nil {
		http.Error(w, "Error getting username bu userid", http.StatusInternalServerError)
		return
	}
	fmt.Printf("getting username from id in session:%s\nGetting user from homepage.user: %s", userName, HomePage.User)

	// Handle image upload
	var imagePath string = ""
	imageFile, header, err := r.FormFile("postImage")
	if err != nil {
		// If there's no image uploaded, leave imagePath empty
		// if err.Error() != "multipart: no mixed boundary parts" {
		// 	log.Printf("Error retrieving image: %v\n", err)
		// 	http.Error(w, "Error uploading image", http.StatusInternalServerError)
		// 	return
		// }
	} else {

		// Validate MIME type (must be an image)
		mimeType := header.Header.Get("Content-Type")
		validMimeTypes := []string{"image/jpeg", "image/png", "image/gif", "image/webp"}

		isValidMimeType := false
		for _, validType := range validMimeTypes {
			if mimeType == validType {
				isValidMimeType = true
				break
			}
		}

		if !isValidMimeType {
			log.Printf("Invalid image type: %s\n", mimeType)
			postCreateError = "Invalid image type. Only JPEG, PNG, GIF, and WEBP are allowed."
			http.Redirect(w, r, "/create-a-post", http.StatusSeeOther)
			fromWhere = "er"
			fmt.Println("error is: ", postCreateError)
			return
			//http.Error(w, "Invalid image type. Only JPEG, PNG, GIF, and WEBP are allowed.", http.StatusBadRequest)
			//return
		}

		const maxFileSize = 10 * 1024 * 1024 // 10 MB

		// Validate file size
		if header.Size > maxFileSize {
			log.Printf("File size exceeds the limit: %d bytes\n", header.Size)
			postCreateError = "File size exceeds the maximum limit of 10MB"
			fromWhere = "er"
			http.Redirect(w, r, "/create-a-post", http.StatusSeeOther)
			return
			// http.Error(w, "File size exceeds the maximum limit of 10MB", http.StatusBadRequest)
			// return
		}

		//check if uploads exist, if not create it
		if _, err := os.Stat("uploads"); os.IsNotExist(err) {
			err := os.Mkdir("uploads", 0755) // Create the uploads directory if it doesn't exist
			if err != nil {
				log.Printf("Error creating uploads directory: %v\n", err)
				http.Error(w, "Error creating uploads directory", http.StatusInternalServerError)
				return
			}
		}

		// Process the image
		fileName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), r.FormValue("postImage"))
		imagePath = fmt.Sprintf("uploads/%s", fileName)
		fmt.Println(imagePath)
		// Save the image to the server
		out, err := os.Create(imagePath)
		if err != nil {
			log.Printf("Failed to save image: %v\n", err)
			http.Error(w, "Error saving image", http.StatusInternalServerError)
			return
		}
		defer out.Close()
		// Copy the uploaded image to the created file
		_, err = io.Copy(out, imageFile)
		if err != nil {
			log.Printf("Failed to copy image: %v\n", err)
			http.Error(w, "Error saving image", http.StatusInternalServerError)
			return
		}
	}

	fmt.Printf("Inserting into posts: title=%s, content=%s, user_id=%d, user_name=%s, category=%s, postTime=%s, image_path=%s\n", ttl, cont, userid, userName, cat, timee, imagePath)

	_, err = db.Exec(`
	INSERT INTO posts 
	(title, content, user_id, user_name, category, postTime, image_path) 
	VALUES (?, ?, ?, ?, ?, ?, ?)`, ttl, cont, userid, userName, cat, timee, imagePath)
	if err != nil {
		log.Printf("Failed to create post: %v\n", err) // Log the error for debugging
		http.Error(w, "Error creating post", http.StatusInternalServerError)
		return
	}
	var maxID int
	if err := db.QueryRow("SELECT MAX(id) FROM posts").Scan(&maxID); err != nil {
		log.Fatal(err)
	}
	var newpost = postsForHome{
		Username: userName,
		Title:    ttl,
		PostId:   maxID,
	}
	HomePage.AllPosts = append([]postsForHome{newpost}, HomePage.AllPosts...)
	postCreateError = ""
	http.Redirect(w, r, "/posted/"+strconv.Itoa(maxID), http.StatusSeeOther)
}

func viewPostHandler(w http.ResponseWriter, r *http.Request) {
	postIdStr := r.URL.Path[len("/posted/"):] // Gets the part of the URL after "/posted/"
	postId, err := strconv.Atoi(postIdStr)    // Convert the ID to an integer
	if err != nil {
		http.Error(w, "Invalid Post ID", http.StatusBadRequest)
		return
	}

	// Fetch the post from the database using the post ID (pseudo-code)
	post, err := fetchOnePost(postId)
	if err != nil {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}
	comments, err := fetchCommentbyPostId(postId)
	//	fmt.Println(comments)
	if err != nil {
		http.Error(w, "error fetching comments", http.StatusBadRequest)
		return
	}
	lili := allPage{
		Cmnts:   *comments,
		Thepost: *post,
	}
	err = templates.ExecuteTemplate(w, "post.html", lili)
	if err != nil {
		http.Error(w, "Unable to render template", http.StatusInternalServerError)
	}
}

func fetchCommentbyPostId(postId int) (*[]myComment, error) {
	query := `
	SELECT u.username, c.content
	FROM comments c
	JOIN users u ON c.user_id = u.id
	WHERE c.post_id = ?;
	`
	rows, err := db.Query(query, postId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var allcomments []myComment
	for rows.Next() {
		var comment myComment
		err := rows.Scan(&comment.Username, &comment.Content)
		if err != nil {
			return nil, err
		}
		allcomments = append(allcomments, comment)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return &allcomments, nil
}

func fetchOnePost(post_id int) (*PostData, error) {
	query := "SELECT id, user_name, content, category, postTime, title, image_path FROM posts WHERE id = ?"
	row := db.QueryRow(query, post_id)

	var postView PostData
	err := row.Scan(&postView.PostID, &postView.Username, &postView.Content, &postView.Categ, &postView.Posttime, &postView.Title, &postView.ImagePath)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no post found with id %d", post_id)
		}
		return nil, err
	}
	fmt.Println(postView)
	postView.LikeCount, err = getLikeCount(postView.PostID)
	if err != nil {
		return nil, err
	}
	postView.DisLikeCount, err = getDislikeCount(postView.PostID)
	if err != nil {
		return nil, err
	}
	postView.CommentCount, err = getCommentCount(postView.PostID)
	if err != nil {
		return nil, err
	}
	return &postView, nil
}
