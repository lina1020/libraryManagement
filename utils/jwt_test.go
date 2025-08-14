package utils

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestGenerateAndParseToken_Success(t *testing.T) {
	userID := uint(123)
	role := "admin"

	token, err := GenerateToken(userID, role)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	claims, err := ParseToken(token)
	assert.NoError(t, err)
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, role, claims.Role)
	assert.WithinDuration(t, time.Now(), claims.RegisteredClaims.IssuedAt.Time, time.Second)
	assert.WithinDuration(t, time.Now().Add(24*time.Hour), claims.RegisteredClaims.ExpiresAt.Time, time.Second)
}

func TestParseToken_ExpiredToken(t *testing.T) {
	// 手动生成一个已过期的 token
	expiredTime := time.Now().Add(-2 * time.Hour)
	claims := &Claims{
		UserID: 1,
		Role:   "user",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiredTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, _ := token.SignedString(jwtSecret)

	_, err := ParseToken(signedToken)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "token has")
}

func TestParseToken_InvalidSignature(t *testing.T) {
	// 使用错误密钥签名的 token
	badSecret := []byte("wrong-secret-key")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &Claims{UserID: 1, Role: "user"})
	signedToken, _ := token.SignedString(badSecret)

	_, err := ParseToken(signedToken)

	assert.Error(t, err)
}
