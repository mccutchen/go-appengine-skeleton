package main

import (
	"time"
)

// User is a user
type User struct {
	Email       string
	Confirmed   bool
	Deactivated bool

	Created time.Time
	Updated time.Time
}
