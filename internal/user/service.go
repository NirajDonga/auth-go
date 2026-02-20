package user

import (
	"context"
	"errors"
	"fmt"
	"go-auth/internal/auth"
	"strings"

	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo      *Repo
	jwtSecret string
}

func NewService(repo *Repo, jwtSecret string) *Service {
	return &Service{
		repo:      repo,
		jwtSecret: jwtSecret,
	}
}

type RegisterInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type Authresult struct {
	Token string     `json:"token"`
	User  PublicUser `json:"user"`
}

func (s *Service) Register(ctx context.Context, input RegisterInput) (Authresult, error) {
	email := strings.ToLower(strings.TrimSpace(input.Email))
	pass := strings.TrimSpace(input.Password)

	if email == "" || pass == "" {
		return Authresult{}, errors.New("Email and Password are required")
	}

	if len(pass) < 6 {
		return Authresult{}, errors.New("Password must be at least 6 characters long")
	}

	_, err := s.repo.FindByEmail(ctx, email)
	if err == nil {
		return Authresult{}, errors.New("Email already registered")
	}

	if !errors.Is(err, pgx.ErrNoRows) {
		return Authresult{}, err
	}

	hashBytes, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		return Authresult{}, fmt.Errorf("Hashing of password failed: %w", err)
	}

	u := User{
		Email:    email,
		Password: string(hashBytes),
		Role:     "user",
	}

	created, err := s.repo.Create(ctx, u)
	if err != nil {
		return Authresult{}, err
	}

	token, err := auth.CreateToken(s.jwtSecret, created.ID, u.Role)
	if err != nil {
		return Authresult{}, err
	}
	return Authresult{
		Token: token,
		User:  ToPublic(created),
	}, nil
}

func (s *Service) LogIn(ctx context.Context, input LoginInput) (Authresult, error) {
	email := strings.ToLower(strings.TrimSpace(input.Email))
	pass := strings.TrimSpace(input.Password)

	if email == "" || pass == "" {
		return Authresult{}, errors.New("Email and Password are required")
	}

	u, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Authresult{}, errors.New("Invalid Credentials")
		}
		return Authresult{}, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(pass)); err != nil {
		return Authresult{}, errors.New("Invalid Credentials or wrong password")
	}

	token, err := auth.CreateToken(s.jwtSecret, u.ID, u.Role)
	if err != nil {
		return Authresult{}, err
	}

	return Authresult{
		Token: token,
		User:  ToPublic(u),
	}, nil

}
