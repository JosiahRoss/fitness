package sessions

import (
	"database/sql"
	"fmt"
	"time"
)

// Database defines the session  database
type Database struct {
	db *sql.DB
}

// New creates a new session  database.
func New(db *sql.DB) *Database {
	return &Database{
		db: db,
	}
}

// Session defines a session.
type Session struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	CreatedOn time.Time `json:"created_on"`
}

// Sessions defines a set of sessions.
type Sessions struct {
	Sessions []*Session `json:"sessions"`
	Total    int        `json:"total"`
}

const (
	// stmtInsert defines teh SQL statement to
	// insert a new session int the database
	stmtInsert = `
INSERT INTO sessions (user_id, created_on)
VALUES(?,?)
`

	// 	// stmtSelect defines the SQL statment to
	// 	// select a set of sessions from sessions.
	// 	stmtSelect = `
	// SELECT id, user_id, created_on
	// FROM sessions
	// %s
	// LIMIT %v, %v
	// `

	// stmtSelectCount defines the SQL statement to
	// select the total number of sessions found for a
	// a given user, according to the filters.
	stmtSelectCount = `
SELECT COUNT(*)
FROM session 
%s
`

	// stmtSelectByID defines the SQL statement to
	// select a Session by its ID.
	stmtSelectByID = `
SELECT id, user_id, created_on
FROM sessions
WHERE id=?
`

	// stmtSelectByUserID defines the SQL statement to
	// select a Session by its ID.
	stmtSelectByUserID = `
SELECT id, user_id, created_on
FROM sessions
WHERE user_id=?
`
	// stmtSelectByIDAndUserID defines the SQL statement to
	// selcet a Session by its id and user_id.
	stmtSelectByIDAndUserID = `
SELECT id, user_id, created_on
FROM sessions
WHERE id=?
	`

	// stmtSelectByUserID defines the SQL statement to
	// select a Session by its ID.
	stmtSelectByCreatedOn = `
SELECT id, user_id, created_on
FROM sessions
WHERE created_on=?
`

	// stmtDelete defines the SQL statement to
	// delete a session
	stmtDelete = `
DELETE FROM sessions
WHERE id=?
AND user_id=?
`
)

// NewParams defines the parameters for the new mehtod.
type NewParams struct {
}

// New creates a and returns a new session.
func (db *Database) New(uid int) (*Session, error) {
	//Create a new Session.
	session := &Session{
		UserID:    uid,
		CreatedOn: time.Now(),
	}

	// Create a variable to hold the result.
	var res sql.Result
	var err error

	// Execute the query.
	if res, err = db.db.Exec(stmtInsert, session.UserID, session.CreatedOn); err != nil {
		return nil, err
	}

	// Get last insert ID.
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	// Define session.ID
	session.ID = int(id)

	return session, nil
}

// GetByID retrieves a session by its ID
func (db *Database) GetByID(id int) (*Session, error) {
	// Create a new Session.
	session := &Session{}

	// Execute the query.
	err := db.db.QueryRow(stmtSelectByID, id).Scan(&session.ID, &session.UserID, &session.CreatedOn)
	switch {
	case err == sql.ErrNoRows:
		return nil, ErrSessionNotFound
	case err != nil:
		return nil, err

	}
	return session, nil

}

// GetByUserID retrieves a session by its userID.
func (db *Database) GetByUserID(uid int) (*Sessions, error) {
	// Create a new Session.
	sessions := &Sessions{
		Sessions: []*Session{},
	}
	// Execute the query.
	rows, err := db.db.Query(stmtSelectByIDAndUserID, uid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Loop though the sessions rows.
	for rows.Next() {
		// Create a new Session.
		session := &Session{}

		// Scan rows values into session struct.
		if err := rows.Scan(&session.ID, &session.UserID, &session.CreatedOn); err != nil {
			return nil, err
		}

		// Add to sessions set
		sessions.Sessions = append(sessions.Sessions, session)

	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	//
	// 	i have worries about the queryCount it feels like it wont work

	//
	//
	//

	field := "WHERE created_on=?"
	// Build the total count query.
	queryCount := fmt.Sprintf(stmtSelectCount, field)

	// Get total count.
	var total int
	if err = db.db.QueryRow(queryCount).Scan(&total); err != nil {
		return nil, err
	}

	sessions.Total = total

	return sessions, nil

}

// GetByIDAndUserID retrieves a set of sessions by id and user_id.
func (db *Database) GetByIDAndUserID(id, uid int) (*Sessions, error) {
	// Create a new slice of sessions.
	sessions := &Sessions{
		Sessions: []*Session{},
	}
	// Execute the query.
	rows, err := db.db.Query(stmtSelectByIDAndUserID, id, uid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Loop though the sessions rows.
	for rows.Next() {
		// Create a new Session.
		session := &Session{}

		// Scan rows values into session struct.
		if err := rows.Scan(&session.ID, &session.UserID, &session.CreatedOn); err != nil {
			return nil, err
		}

		// Add to sessions set
		sessions.Sessions = append(sessions.Sessions, session)

	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	field := "WHERE created_on=?"
	// Build the total count query.
	queryCount := fmt.Sprintf(stmtSelectCount, field)

	// Get total count.
	var total int
	if err = db.db.QueryRow(queryCount).Scan(&total); err != nil {
		return nil, err
	}

	sessions.Total = total

	return sessions, nil

}

// GetByCreatedOn retrieves a group of session by
// the created_on field.
func (db *Database) GetByCreatedOn(date time.Time) (*Sessions, error) {
	// Create a new slice of sessions
	sessions := &Sessions{
		Sessions: []*Session{},
	}
	// Execute the query.
	rows, err := db.db.Query(stmtSelectByCreatedOn, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Loop though the sessions rows.
	for rows.Next() {
		// Create a new Session.
		session := &Session{}

		// Scan rows values into session struct.
		if err := rows.Scan(&session.ID, &session.UserID, &session.CreatedOn); err != nil {
			return nil, err
		}

		// Add to sessions set
		sessions.Sessions = append(sessions.Sessions, session)

	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	field := "WHERE created_on=?"
	// Build the total count query.
	queryCount := fmt.Sprintf(stmtSelectCount, field)

	// Get total count.
	var total int
	if err = db.db.QueryRow(queryCount).Scan(&total); err != nil {
		return nil, err
	}

	sessions.Total = total

	return sessions, nil

}
