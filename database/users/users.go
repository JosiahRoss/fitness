package users

import (
	"database/sql"
	"time"
)

// Database defines the members database.
type Database struct {
	db *sql.DB
}

// New creates a new User database.
func New(db *sql.DB) *Database {
	return &Database{
		db: db,
	}
}

// User defines a user.
type User struct {
	ID        int       `json:"id"`
	FullName  string    `json:"full_name"`
	UserName  string    `json:"user_name"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
}

const (
	// stmtInsert defines the SQL statement to
	// insert a new user into the database.
	stmtInsert = `
INSERT INTO users (full_name, user_name,email, password, created_at)
VALUES (?, ?, ?, ?, ?)
`

	// stmtSelectByID defines the SQL statement to
	// select a user by their ID.
	stmtSelectByID = `
SELECT id,full_name, user_name, email, password created_at
FROM users
WHERE id=?
`

	// stmtSelectByEmail defines the SQL statement
	// to select a user by their email address.
	stmtSelectByEmail = `
SELECT id,full_name, user_name, email, password created_at
FROM users
WHERE email=?
`
)

// NewParams defines the parameter for the New Method.
type NewParams struct {
	FullName  string    `json:"full_name"`
	UserName  string    `json:"user_name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
}

// New creates a new use.
func (db *Database) New(params *NewParams) (*User, error) {
	// Create a new user.
	user := &User{
		FullName:  params.FullName,
		UserName:  params.UserName,
		Email:     params.Email,
		Password:  params.Password,
		CreatedAt: params.CreatedAt,
	}

	// Create variable to hold the result.
	var res sql.Result
	var err error

	// Execute the query.
	if res, err = db.db.Exec(stmtInsert, user.FullName, user.UserName, user.Email, user.Password, user.CreatedAt); err != nil {
		return nil, err
	}

	// Get last insert ID.
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	user.ID = int(id)

	return user, nil
}

// GetByID retrieves a user by their ID.
func (db *Database) GetByID(id int) (*User, error) {
	// Create a new User.
	user := &User{}

	// Execute the query.
	err := db.db.QueryRow(stmtSelectByID, id).Scan(&user.ID, user.FullName, user.UserName, user.Email, user.Password, user.CreatedAt)
	switch {
	case err == sql.ErrNoRows:
		return nil, ErrUserNotFound
	case err != nil:
		return nil, err
	}

	return user, nil
}

// GetByEmail retrieves a user by their email.
func (db *Database) GetByEmail(email string) (*User, error) {
	// Create a new Member.
	user := &User{}

	// Execute the query.
	err := db.db.QueryRow(stmtSelectByEmail, email).Scan(&user.ID, user.FullName, user.UserName, user.Email, user.Password, user.CreatedAt)
	switch {
	case err == sql.ErrNoRows:
		return nil, ErrUserNotFound
	case err != nil:
		return nil, err
	}

	return user, nil
}
