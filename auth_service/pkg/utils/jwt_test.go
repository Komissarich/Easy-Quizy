package utils

import (
	"fmt"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateJWTAndParseJWT(t *testing.T) {
	secret := "test-secret"
	userID := uint64(123)
	email := "test@example.com"

	tokenString, err := GenerateJWT(secret, userID, email)
	require.NoError(t, err)
	require.NotEmpty(t, tokenString)

	claims, err := ParseJWT(tokenString, secret)
	require.NoError(t, err)

	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, email, claims.Email)
	assert.True(t, claims.ExpiresAt.After(time.Now()))
}

func TestParseJWT_InvalidToken(t *testing.T) {
	testCases := []struct {
		name        string
		tokenString string
		secret      string
		expectedErr string
	}{
		{
			name:        "empty token",
			tokenString: "",
			secret:      "secret",
			expectedErr: "empty token string",
		},
		{
			name:        "invalid signing method",
			tokenString: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxMjMsImVtYWlsIjoidGVzdEBleGFtcGxlLmNvbSJ9.7QJwC0Zq7X9X7Z3X2X1X0X9X8X7X6X5X4X3X2X1X0",
			secret:      "secret",
			expectedErr: "unexpected signing method",
		},
		{
			name:        "invalid secret",
			tokenString: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxMjMsImVtYWlsIjoidGVzdEBleGFtcGxlLmNvbSJ9.7QJwC0Zq7X9X7Z3X2X1X0X9X8X7X6X5X4X3X2X1X0",
			secret:      "wrong-secret",
			expectedErr: "failed to parse token",
		},
		{
			name:        "malformed token",
			tokenString: "malformed.token.string",
			secret:      "secret",
			expectedErr: "failed to parse token",
		},
		{
			name: "expired token",
			tokenString: func() string {
				claims := Claims{
					UserID: 123,
					Email:  "test@example.com",
					RegisteredClaims: jwt.RegisteredClaims{
						ExpiresAt: jwt.NewNumericDate(time.Now().Add(-24 * time.Hour)),
					},
				}
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
				tokenString, _ := token.SignedString([]byte("secret"))
				return tokenString
			}(),
			secret:      "secret",
			expectedErr: "invalid token claims",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := ParseJWT(tc.tokenString, tc.secret)
			require.Error(t, err)
			assert.Contains(t, err.Error(), tc.expectedErr)
		})
	}
}

func TestGenerateJWT_EmptySecret(t *testing.T) {
	_, err := GenerateJWT("", 123, "test@example.com")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "empty secret")
}

func TestExtractTokenFromHeader(t *testing.T) {
	testCases := []struct {
		name        string
		header      string
		expected    string
		expectedErr string
	}{
		{
			name:        "valid header",
			header:      "Bearer valid.token.string",
			expected:    "valid.token.string",
			expectedErr: "",
		},
		{
			name:        "valid header with spaces",
			header:      "   Bearer    valid.token.string   ",
			expected:    "valid.token.string",
			expectedErr: "",
		},
		{
			name:        "empty header",
			header:      "",
			expected:    "",
			expectedErr: "empty authorization header",
		},
		{
			name:        "invalid scheme",
			header:      "Basic base64string",
			expected:    "",
			expectedErr: "invalid auth scheme, expected Bearer",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			token, err := ExtractTokenFromHeader(tc.header)

			if tc.expectedErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedErr)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expected, token)
			}
		})
	}
}

func TestIntegration_GenerateParseExtract(t *testing.T) {
	secret := "integration-secret"
	userID := uint64(42)
	email := "integration@test.com"

	token, err := GenerateJWT(secret, userID, email)
	require.NoError(t, err)

	header := fmt.Sprintf("Bearer %s", token)

	extractedToken, err := ExtractTokenFromHeader(header)
	require.NoError(t, err)
	assert.Equal(t, token, extractedToken)

	claims, err := ParseJWT(extractedToken, secret)
	require.NoError(t, err)

	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, email, claims.Email)
}
