package service

import "errors"

var (
	ErrUserNotFound        = errors.New("user not found")
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrCannotAddYourself   = errors.New("cannot add yourself as a friend")
	ErrAlreadyFriends      = errors.New("users are already friends")
	ErrNotFriends          = errors.New("users are not friends")
	ErrFriendNotFound      = errors.New("friend not found")
	ErrOperationNotAllowed = errors.New("operation not allowed")
)
