package user

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repo struct {
	pool *pgxpool.Pool
}

func NewRepo(pool *pgxpool.Pool) *Repo {
	return &Repo{pool: pool}
}

func (r *Repo) FindByEmail(ctx context.Context, email string) (User, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	childCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var u User
	row := r.pool.QueryRow(childCtx, `SELECT id, email, password, role, created_at, updated_at FROM users WHERE email=$1`, email)
	if err := row.Scan(&u.ID, &u.Email, &u.Password, &u.Role, &u.CreatedAt, &u.UpdatedAt); err != nil {
		if err == pgx.ErrNoRows {
			return User{}, pgx.ErrNoRows
		}
		return User{}, fmt.Errorf("find by email failed: %w", err)
	}
	return u, nil
}

func (r *Repo) Create(ctx context.Context, u User) (User, error) {
	id := uuid.NewString()
	now := time.Now().UTC()
	u.ID = id
	u.CreatedAt = now
	u.UpdatedAt = now

	childCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := r.pool.Exec(childCtx, `INSERT INTO users (id, email, password, role, created_at, updated_at) VALUES ($1,$2,$3,$4,$5,$6)`,
		u.ID, u.Email, u.Password, u.Role, u.CreatedAt, u.UpdatedAt)
	if err != nil {
		return User{}, fmt.Errorf("insert user failed: %w", err)
	}

	return u, nil
}
