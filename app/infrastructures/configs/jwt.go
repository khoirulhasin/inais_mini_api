package configs

import (
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/khoirulhasin/untirta_api/app/infrastructures/constants"
)

func GetClaim(userId int) jwt.RegisteredClaims {
	return jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(constants.ExpiresAt * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		NotBefore: jwt.NewNumericDate(time.Now()),
		Issuer:    constants.Issuer,
		Subject:   constants.Subject,
		// convert userId to string
		ID:       strconv.Itoa(userId),
		Audience: []string{"our_client"},
	}
}
