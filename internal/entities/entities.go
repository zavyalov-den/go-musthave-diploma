package entities

import (
	"errors"
	"time"
)

var ErrUserConflict = errors.New("requested db entry is created by different user")
var ErrNoContent = errors.New("no content")

type Credentials struct {
	UserID   int
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Order struct {
	Number     int       `json:"number"`
	Status     string    `json:"status"`
	Accrual    int       `json:"accrual,omitempty"`
	UploadedAt time.Time `json:"uploaded_at"`
}
