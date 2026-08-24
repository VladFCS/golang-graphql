package user

import "errors"

var ErrNotFound = errors.New("user not found")
var ErrDuplicateUser = errors.New("user already exists")
var ErrEmailMustPresent = errors.New("email must exist")
var ErrUsernameMustPresent = errors.New("username must exist")
var ErrPasswordMustPresent = errors.New("password must exist")