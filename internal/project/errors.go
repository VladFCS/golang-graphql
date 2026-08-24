package project

import "errors"

var ErrNotFound = errors.New("project not found")
var ErrNameMustPresent = errors.New("project name must exist")
var ErrCreatedByMustPresent = errors.New("project creator must exist")
