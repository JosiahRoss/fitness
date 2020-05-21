package users

import (
	"fitness/database"
	"fitness/services/errors"
	"fmt"
	"time"

	dbusers "fitness/database/users"

	"golang.org/x/crypto/bcrypt"
)

// Service defines the users service.
type Service struct {
	db *database.Database
}

// New returns a new users service.
func New(db *database.Database) *Service {
	return &Service{
		db: db,
	}
}

// User defines a user.
type User dbusers.User

// NewParams defines the parameters for the New method.
type NewParams dbusers.NewParams

// New creates a new user.
func (s *Service) New(params *NewParams) (*User, error) {
	//fmt.Println("service.New()")
	// Create a new ParamErrors.
	pes := errors.NewParamErrors()

	// Check full_name.
	if params.FullName == "" {

		pes.Add(errors.NewParamError("full_name", ErrFullNameEmpty))
	}
	//fmt.Println("checking full name")

	// Check user_name.
	if params.UserName == "" {
		pes.Add(errors.NewParamError("user_name", ErrUserNameEmpty))
	} else {
		_, err := s.db.Users.GetByUserName(params.UserName)
		if err == nil {
			fmt.Println(err)
			pes.Add(errors.NewParamError("user_name", ErrUserNameExists))
		} else if err != nil && err != dbusers.ErrUserNotFound {
			fmt.Println(err)
			return nil, err
		}
	}
	fmt.Println("checking user name")

	// Check email.
	if params.Email == "" {
		pes.Add(errors.NewParamError("email", ErrEmailEmpty))
	} else {
		_, err := s.db.Users.GetByEmail(params.Email)
		if err == nil {
			pes.Add(errors.NewParamError("email", ErrEmailExists))
		} else if err != nil && err != dbusers.ErrUserNotFound {
			return nil, err
		}
	}
	fmt.Println("checking email")

	// Check password.
	if len(params.Password) < 8 {
		pes.Add(errors.NewParamError("password", ErrPassword))
	}
	fmt.Println("checking password")

	// Return if there were parameter errors.
	if pes.Length() > 0 {
		return nil, pes
	}

	// Hash the password.
	pwhash, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create this user in the database.
	dbu, err := s.db.Users.New(&dbusers.NewParams{
		FullName:  params.FullName,
		UserName:  params.UserName,
		Email:     params.Email,
		Password:  string(pwhash),
		CreatedAt: time.Now(),
	})
	if err != nil {
		return nil, err
	}

	// Create a new user.
	user := &User{
		ID:        dbu.ID,
		FullName:  dbu.FullName,
		UserName:  dbu.UserName,
		Email:     dbu.Email,
		Password:  dbu.Password,
		CreatedAt: dbu.CreatedAt,
	}

	return user, nil
}

//LoginParams defines tehparameters for the Login method.
type LoginParams struct {
	Email    string
	Password string
}

// Login checks if a User exists in the database and can log in.
func (s *Service) Login(params *LoginParams) (*User, error) {
	// Try t o pull this user from the database.
	dbu, err := s.db.Users.GetByEmail(params.Email)
	//fmt.Println(err)
	if err == dbusers.ErrUserNotFound {
		return nil, ErrInvalidLogin
	} else if err != nil {
		return nil, err
	}

	// Validate the password.
	if err = bcrypt.CompareHashAndPassword([]byte(dbu.Password), []byte(params.Password)); err != nil {
		return nil, ErrInvalidLogin
	}

	// Create a new User.
	user := &User{
		ID:        dbu.ID,
		FullName:  dbu.FullName,
		UserName:  dbu.UserName,
		Email:     dbu.Email,
		Password:  dbu.Password,
		CreatedAt: dbu.CreatedAt,
	}

	return user, nil

}

// GetByID retrieves a User by their ID.
func (s *Service) GetByID(id int) (*User, error) {
	// Try to pull this User from the database.
	dbm, err := s.db.Users.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Create a new User.
	user := &User{
		ID:       dbm.ID,
		Email:    dbm.Email,
		Password: dbm.Password,
	}

	return user, nil
}
