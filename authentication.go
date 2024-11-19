package main

import (
	"regexp"
	"strings"
	"unicode"
)

// Password validation logic
func checkpassword(s string) (bool, string) {
	var upper, lower, num, special bool

	if len(s) < 8 {
		return false, ("Password is less than 8 characters\n")
	}
	if strings.Contains(s, " ") {
		return false, ("Password cannot contain spaces\n")
	}
	for _, i := range s {
		if unicode.IsUpper(i) {
			upper = true
		} else if unicode.IsDigit(i) {
			num = true
		} else if unicode.IsLower(i) {
			lower = true
		} else if unicode.IsSymbol(i) || unicode.IsPunct(i) {
			special = true
		}
		if upper && lower && num && special {
			return true, ""
		}
	}
	lolo := ""
	if !upper {
		lolo += ("Missing uppercase letter\n")
	}
	if !lower {
		lolo += ("Missing lowercase letter\n")
	}
	if !num {
		lolo += ("Missing number\n")
	}
	if !special {
		lolo += ("Missing special character\n")
	}
	return (upper && lower && num && special), lolo
}

// Email validation logic
func checkEmail(email string) bool {
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.com$`
	re := regexp.MustCompile(emailRegex)
	return re.MatchString(email)
}

func getLikeCount(postID int) (int, error) {
	var likeCount int
	err := db.QueryRow("SELECT COUNT(*) FROM likes WHERE post_id = ? AND is_like = true", postID).Scan(&likeCount)
	if err != nil {
		return 0, err
	}
	return likeCount, nil
}

func getDislikeCount(postID int) (int, error) {
	var dislikeCount int
	err := db.QueryRow("SELECT COUNT(*) FROM likes WHERE post_id = ? AND is_like = false", postID).Scan(&dislikeCount)
	if err != nil {
		return 0, err
	}
	return dislikeCount, nil
}
