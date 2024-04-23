package auth

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

const Expiration = time.Hour

var jwtKey = []byte("my_secret_key")

// Claims структура, включает стандартные jwt.Claims и пользовательские поля
type Claims struct {
	IsAdmin bool `json:"is_admin"`
	jwt.StandardClaims
}

type Auth struct {
}

func NewAuth() *Auth {
	return &Auth{}
}

func (a *Auth) Authenticate(token string) error {
	_, err := validateJWT(token)
	if err != nil {
		return err
	}
	return nil
}

func (a *Auth) IsAdmin(token string) bool {
	claims, _ := validateJWT(token)
	if claims.IsAdmin {
		return true
	}
	return false
}

func (a *Auth) GenerateToken(isAdmin bool) (string, error) {
	expirationTime := time.Now().Add(1 * Expiration)
	claims := &Claims{
		IsAdmin: isAdmin,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func validateJWT(token string) (*Claims, error) {
	jwtToken, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := jwtToken.Claims.(*Claims); ok && jwtToken.Valid {
		return claims, nil
	} else {
		return nil, fmt.Errorf("invalid JWT token")
	}
}
