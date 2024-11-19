package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

var (
	HomePage        *data
	db              *sql.DB
	AllHomePosts    []postsForHome
	signInErr       FormError
	templates       *template.Template
	postCreateError string
	fromWhere       string
)

// Initialize and parse templates
func init() {
	var err error
	templates, err = template.ParseFiles("./templates/myhome.html", "./templates/create.html", "./templates/index.html", "./templates/post.html", "./templates/sign-in.html")
	if err != nil {
		log.Fatal("Error parsing templates: ", err)
	}
}
func generate_LoginPosts() ([]postsForHome, error) {
	var err error
	var maxID int
	if err = db.QueryRow("SELECT MAX(id) FROM posts").Scan(&maxID); err != nil {
		log.Fatal(err)
	}

	guest_posts := []postsForHome{}
	if err != nil {
		return nil, err
	}
	for i := 0; i < 5; i++ {
		randomID := rand.Intn(maxID)
		if randomID < len(HomePage.AllPosts) {
			guest_posts = append(guest_posts, HomePage.AllPosts[randomID])
		}
	}
	return guest_posts, nil
}

func fetchHomePosts() ([]postsForHome, error) {
	var homepageData []postsForHome
	query := `
		SELECT id, user_name, title
		FROM posts 
		ORDER BY id DESC;
	`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var post postsForHome
		if err := rows.Scan(&post.PostId, &post.Username, &post.Title); err != nil {
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
		homepageData = append(homepageData, post)
	}

	return homepageData, nil
}

func main() {
	var err error
	db, err = initDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	AllHomePosts, err = fetchHomePosts()
	if err != nil {
		fmt.Println("error fetching home page's posts ", err)
		return
	}
	HomePage = &data{}

	HomePage.AllPosts = AllHomePosts

	displayPosts, err := generate_LoginPosts()
	if err != nil {
		fmt.Println("error getting index page's posts ", err)
		return
	}
	http.HandleFunc("/", homeGuest(displayPosts)) // index page guest

	http.HandleFunc("/sign-in", loginPage(&signInErr)) // log-in btn

	http.HandleFunc("/register", RegisterNewUser) // sign guest then to home
	http.HandleFunc("/login", UserLogin)          // find guest then home

	http.HandleFunc("/home", homePageTmpl) // homepage

	http.HandleFunc("/create-a-post", creatApost) // enter post data
	http.HandleFunc("/posted", postPostTodb)      // push post to db
	http.HandleFunc("/posted/", viewPostHandler)  // This handles /posted/{id} URLs

	http.HandleFunc("/submit-comment/", commentSubmit) // push comment

	http.HandleFunc("/log-out", logOutHandler)
	http.HandleFunc("/category/", filterPostsHandler)
	http.HandleFunc("/like", likePost)
	http.HandleFunc("/my-posts", mypostsHandler)

	http.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("./uploads"))))
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./templates/css"))))

	fmt.Println("Listening on port 8000")
	http.ListenAndServe(":8000", nil)
}
