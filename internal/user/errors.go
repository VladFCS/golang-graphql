package user

import "errors"

var ErrNotFound = errors.New("user not found")
var ErrDuplicateUser = errors.New("user already exists")