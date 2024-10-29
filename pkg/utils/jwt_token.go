package utils

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/paudelgaurav/gin-boilerplate/pkg/constants"
)

// GenerateTokens generates access and refresh tokens for a given user's email)
func GenerateTokens(jwtSecret, email string) (string, string, error) {
	// Access token claims
	accessTokenClaims := jwt.MapClaims{}
	accessTokenClaims["authorized"] = true
	accessTokenClaims["email"] = email
	accessTokenClaims["exp"] = time.Now().Add(constants.AcessTokenExpiry).Unix() // 15 mins expiry
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)

	secret := []byte(jwtSecret)

	// Sign the access token
	signedAccessToken, err := accessToken.SignedString(secret)
	if err != nil {
		return "", "", err
	}

	// Refresh token claims
	refreshTokenClaims := jwt.MapClaims{}
	refreshTokenClaims["email"] = email
	refreshTokenClaims["exp"] = time.Now().Add(constants.RefreshTokenExpiry).Unix() // 3 days expiry
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)

	// Sign the refresh token
	signedRefreshToken, err := refreshToken.SignedString(secret)
	if err != nil {
		return "", "", err
	}

	return signedAccessToken, signedRefreshToken, nil
}
