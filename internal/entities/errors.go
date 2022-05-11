package entities

import "errors"

var ErrUserConflict = errors.New("requested db entry is created by different user")
var ErrEntryExists = errors.New("db entry is already crated")
var ErrNoContent = errors.New("no content")
