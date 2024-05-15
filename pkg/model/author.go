package model

import (
	"fmt"
	"os"
	"time"
)

type Author struct {
	Name  string
	Email string
	T     time.Time
}

func (a *Author) toStr() string {
	authorStr := fmt.Sprintf("%s %s %s", a.Name, a.Email, a.T.String())
	return authorStr
}

func NewAuthor() *Author {

	// Fetching Author Name and email from Environment variable, you can set them in your ~/.profile
	name := os.Getenv("SASY_AUTHOR_NAME")
	email := os.Getenv("SASY_AUTHOR_EMAIL")
	author := &Author{Name: name, Email: email, T: time.Now()}
	return author
}
