package entities

import (
	"errors"
	"time"
)

var ErrUserConflict = errors.New("requested db entry is created by different user")
var ErrEntryExists = errors.New("db entry is already crated")
var ErrNoContent = errors.New("no content")

type Credentials struct {
	UserID   int
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Order struct {
	Number     string    `json:"number"`
	Status     string    `json:"status"`
	Accrual    float32   `json:"accrual,omitempty"`
	UploadedAt time.Time `json:"uploaded_at,omitempty"`
}

type AccrualOrder struct {
	Order      string    `json:"order"`
	Status     string    `json:"status"`
	Accrual    float32   `json:"accrual,omitempty"`
	UploadedAt time.Time `json:"uploaded_at,omitempty"`
}

type Balance struct {
	Current   float32 `json:"current"`
	Withdrawn float32 `json:"withdrawn"`
}
