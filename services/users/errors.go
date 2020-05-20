package users

import (
	"errors"

	dbusers "fitness/database/users"
)

var (
	// ErrEmailEmpty is returned when the email param is empty.
	ErrEmailEmpty = errors.New("Email parameter is empty")

	// ErrEmailExists is returned when the email already exists.
	ErrEmailExists = errors.New("Email already exists")

	// ErrPassword is returned when the password is in an invalid format.
	ErrPassword = errors.New("Password must be at least 8 characters")

	// ErrInvalidLogin is returned when the email and/or password used
	// with login is invalid.
	ErrInvalidLogin = errors.New("Email and/or password is invalid")

	// ErrUserNotFound is returned when a user could not be found.
	ErrUserNotFound = dbusers.ErrUserNotFound

	// ErrFullNameEmpty is returned when the full_name param is empty.
	ErrFullNameEmpty = errors.New("Full Name parameter is empty")

	// ErrUserNameEmpty is returned when the user_name param is empty.
	ErrUserNameEmpty = errors.New("User Name is empty")

	// ErrUserNameExists is returned when the email already exists.
	ErrUserNameExists = errors.New("User Name  already exists")
)
