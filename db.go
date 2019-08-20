package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"time"

	"cloud.google.com/go/datastore"
)

// DB is a layer of business logic on top of an underlying datastore
type DB struct {
	client *datastore.Client
}

// NewDB creates a new DB
func NewDB(client *datastore.Client) *DB {
	return &DB{
		client: client,
	}
}

// GetUser gets a user by email
func (db *DB) GetUser(ctx context.Context, email string) (*User, error) {
	user := &User{}
	err := db.client.Get(ctx, userKey(email), user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// CreateUser creates and stores a new User
func (db *DB) CreateUser(ctx context.Context, email string) (*User, error) {
	now := time.Now()
	user := &User{
		Email:   email,
		Created: now,
		Updated: now,
	}
	_, err := db.client.Put(ctx, userKey(email), user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// PutUser updates an existing User
func (db *DB) PutUser(ctx context.Context, user *User) error {
	user.Updated = time.Now()
	_, err := db.client.Put(ctx, userKey(user.Email), user)
	return err
}

func userKey(email string) *datastore.Key {
	return datastore.NameKey("User", email, nil)
}

func hashStrings(xs ...string) string {
	h := sha256.New()
	return hex.EncodeToString(h.Sum(joinForKey(xs...)))
}

func joinForKey(xs ...string) []byte {
	return []byte(strings.Join(xs, "/"))
}
