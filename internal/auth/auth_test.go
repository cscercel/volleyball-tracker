package auth

import (
	"errors"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestValidateJWT(t *testing.T) {
	userID := uuid.New()
	validToken, _ := MakeJWT(userID, "secret", time.Hour)

	tests := []struct {
		name        string
		tokenString string
		tokenSecret string
		wantUserID  uuid.UUID
		wantErr     bool
	}{
		{
			name:        "Valid token",
			tokenString: validToken,
			tokenSecret: "secret",
			wantUserID:  userID,
			wantErr:     false,
		},
		{
			name:        "Invalid token",
			tokenString: "invalid.token.string",
			tokenSecret: "secret",
			wantUserID:  uuid.Nil,
			wantErr:     true,
		},
		{
			name:        "Wrong secret",
			tokenString: validToken,
			tokenSecret: "wrong_secret",
			wantUserID:  uuid.Nil,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotUserID, err := ValidateJWT(tt.tokenString, tt.tokenSecret)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateJWT() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotUserID != tt.wantUserID {
				t.Errorf("ValidateJWT() gotUserID = %v, want %v", gotUserID, tt.wantUserID)
			}
		})
	}
}

func TestGetBearerToken(t *testing.T) {
	req := httptest.NewRequest("GET", "/path", nil)
	req.Header.Set("Authorization", "Bearer token123")

	invalid_req := httptest.NewRequest("GET", "/path", nil)

	want_token := "token123"
	want_err := errors.New("No Authorization header found")

	// Validate token
	got_token, _ := GetBearerToken(req.Header)
	if got_token != want_token {
		t.Errorf("GetBearerToken() got_token = %v, want_token = %v", got_token, want_token)
	}

	// Validate error
	_, err := GetBearerToken(invalid_req.Header)
	if (err != nil) != true {
		t.Errorf("GetBearerToken() error = %v, want_err = %v", err, want_err)
	}
}
