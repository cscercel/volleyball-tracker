package service

import (
	"context"
	"fmt"
	"time"

	"github.com/cscercel/volleyball-tracker/internal/auth"
	"github.com/cscercel/volleyball-tracker/internal/db"
	"github.com/google/uuid"
)

type UserLogin struct {
	ID    uuid.UUID `json:"id"`
	Email string    `json:"email"`
	Token string    `json:"token"`
}

type UserService struct {
	queries          *db.Queries
	registrationCode string
	jwtSecret        string
}

func NewUserService(
	queries *db.Queries, registration_code, jwt_secret string,
) *UserService {
	return &UserService{
		queries:          queries,
		registrationCode: registration_code,
		jwtSecret:        jwt_secret,
	}
}

func (s *UserService) CreateUser(
	ctx context.Context, email, password, registration_code string,
) (db.CreateUserRow, error) {
	// Check if valid registration code was provided
	if registration_code != s.registrationCode {
		return db.CreateUserRow{}, fmt.Errorf("invalid registration code")
	}

	// Hash password
	hashed_password, err := auth.HashPassword(password)
	if err != nil {
		return db.CreateUserRow{}, fmt.Errorf("could not hash password: %w", err)
	}

	user, err := s.queries.CreateUser(ctx, db.CreateUserParams{
		Email:          email,
		HashedPassword: hashed_password,
	})
	if err != nil {
		return db.CreateUserRow{}, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

func (s *UserService) GetUser(ctx context.Context, email, password string) (UserLogin, error) {
	user, err := s.queries.GetUserByEmail(ctx, email)
	if err != nil {
		return UserLogin{}, fmt.Errorf("incorrect email or password: %w", err)
	}

	// Check if password is valid
	is_valid, err := auth.CheckPasswordHash(password, user.HashedPassword)
	if err != nil {
		return UserLogin{}, fmt.Errorf("incorrect email or password: %w", err)
	}

	if !is_valid {
		return UserLogin{}, fmt.Errorf("incorrect email or password: %w", err)
	}

	// Generate JWT token
	token, err := auth.MakeJWT(user.ID, s.jwtSecret, time.Hour)
	if err != nil {
		return UserLogin{}, fmt.Errorf("failed to generate token: %w", err)
	}

	return UserLogin{
		ID:    user.ID,
		Email: user.Email,
		Token: token,
	}, nil
}

func (s *UserService) UpdateUserEmail(
	ctx context.Context, id uuid.UUID, new_email string,
) (db.UpdateUserEmailRow, error) {
	user, err := s.queries.GetUserByID(ctx, id)
	if err != nil {
		return db.UpdateUserEmailRow{}, fmt.Errorf("failed to retrieve user: %w", err)
	}

	user_updated, err := s.queries.UpdateUserEmail(ctx, db.UpdateUserEmailParams{
		ID:    user.ID,
		Email: new_email,
	})
	if err != nil {
		return db.UpdateUserEmailRow{}, fmt.Errorf("failed to change email: %w", err)
	}

	return user_updated, nil
}

func (s *UserService) UpdateUserPassword(
	ctx context.Context, id uuid.UUID, new_password string,
) (db.UpdateUserPasswordRow, error) {
	user, err := s.queries.GetUserByID(ctx, id)
	if err != nil {
		return db.UpdateUserPasswordRow{}, fmt.Errorf("failed to retrieve user: %w", err)
	}

	// Hash new password
	hashed_password, err := auth.HashPassword(new_password)
	if err != nil {
		return db.UpdateUserPasswordRow{}, fmt.Errorf("could not hash password: %w", err)
	}

	user_updated, err := s.queries.UpdateUserPassword(ctx, db.UpdateUserPasswordParams{
		ID:             user.ID,
		HashedPassword: hashed_password,
	})
	if err != nil {
		return db.UpdateUserPasswordRow{}, fmt.Errorf("failed to change password: %w", err)
	}

	return user_updated, nil
}
