package entities

import "errors"

var ErrUserConflict = errors.New("requested db entry is created by different user")

type Credentials struct {
	UserID   int
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Order struct {
	Number     string `json:"number"`
	Status     string `json:"status"`
	Accrual    int    `json:"accrual,omitempty"`
	UploadedAt string `json:"uploaded_at"`
}
