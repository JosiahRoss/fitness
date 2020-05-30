package sessions

import (
	"fitness/database"
	dbsessions "fitness/database/sessions"
)

// Service defines the sessions service.
type Service struct {
	db *database.Database
}

// New returns a new sessions service.
func New(db *database.Database) *Service {
	return &Service{
		db: db,
	}
}

// Session defines a session.
type Session dbsessions.Session

// Sessions defines a set of sessions.
type Sessions struct {
	Sessions []*Session `json:"sessions"`
	Total    int        `json:"total"`
}

// New creates a new session.
func (s *Service) New(uid int) (*Session, error) {

	// Create a this user in the database
	dbs, err := s.db.Sessions.New(uid)
	if err != nil {
		return nil, err
	}

	// Create a new session.
	session := &Session{
		ID:        dbs.ID,
		UserID:    dbs.UserID,
		CreatedOn: dbs.CreatedOn,
	}

	return session, nil
}

// GetByID retrieves a session by its ID.
func (s *Service) GetByID(id int) (*Session, error) {
	dbs, err := s.db.Sessions.GetByID(id)
	if err != nil {
		return nil, err
	}
	// Create a new Session.
	session := &Session{
		ID:        dbs.ID,
		UserID:    dbs.UserID,
		CreatedOn: dbs.CreatedOn,
	}

	return session, nil
}

// GetByUserID retrieves a set of sessions by its UserID
