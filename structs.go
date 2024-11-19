package main

type FormError struct {
	Error   string
	Success string
}
type data struct {
	User     string
	AllPosts []postsForHome
}
type postCreate struct {
	Username string
	DateNow  string
	Err      string
}

type postsForHome struct {
	Username     string
	Title        string
	PostId       int
	LikeCount    int
	DisLikeCount int
	CommentCount int
}
type PostData struct {
	Username     string
	Content      string
	Categ        string
	Posttime     string
	Title        string
	PostID       int
	ImagePath    string // Add this to store the image path
	LikeCount    int
	DisLikeCount int
	CommentCount int
}
type myComment struct {
	Username string
	Content  string
}
type allPage struct {
	Cmnts   []myComment
	Thepost PostData
}
