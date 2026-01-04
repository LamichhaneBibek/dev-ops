package store

import (
	"context"
	"database/sql"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type UserStore struct {
	db *sqlx.DB
}

func NewUserStore(db *sql.DB) *UserStore {
	return &UserStore{db: sqlx.NewDb(db, "postgres")}
}

type User struct {
	Id             uuid.UUID `db:"id"`
	Email          string    `db:"email"`
	HashedPassword string    `db:"hashed_password"`
	CreatedAt      time.Time `db:"created_at"`
}

func (us *User) ComparePassword(password string) error {
	hashPassword, err := base64.StdEncoding.DecodeString(us.HashedPassword)
	if err != nil {
		return fmt.Errorf("failed to decode hash password: %w", err)
	}

	err = bcrypt.CompareHashAndPassword(hashPassword, []byte(password))
	if err != nil {
		return fmt.Errorf("invalid password: %w", err)
	}
	return nil
}

func (us *UserStore) CreateUser(ctx context.Context, email string, password string) (*User, error) {
	const query = `INSERT INTO users (email, hashed_password) VALUES ($1, $2) RETURNING *`
	var user User
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to generate hash password: %w", err)
	}

	hashedPasswordBase64 := base64.StdEncoding.EncodeToString(bytes)

	err = us.db.GetContext(ctx, &user, query, email, hashedPasswordBase64)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	return &user, nil
}

func (us *UserStore) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	const query = `SELECT * FROM users WHERE email = $1`
	var user User
	err := us.db.GetContext(ctx, &user, query, email)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	return &user, nil
}

func (us *UserStore) GetUserById(ctx context.Context, userId uuid.UUID) (*User, error) {
	const query = `SELECT * FROM users WHERE id = $1`
	var user User
	err := us.db.GetContext(ctx, &user, query, userId)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}
	return &user, nil
}
