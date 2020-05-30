package sessions

import "time"

// Session
type Session struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	CreatedOn time.Time `json:"created_on`
}

// Meta defines the response top level meta object.
type Meta struct {
	Total int `json:"total"`
}

// ResultG
