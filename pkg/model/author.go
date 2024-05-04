package model

import (
	"fmt"
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