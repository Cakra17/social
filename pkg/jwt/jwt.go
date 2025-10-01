package jwt

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTUser struct {
	ID string `json:"userId"`
	Email  string `json:"email"`
}

type JWTAuthenticator struct {
	secret string
	duration time.Duration
}

func NewJWTAuthenticator(secret string) *JWTAuthenticator {
	return &JWTAuthenticator{
		secret: secret,
	}
}

func (ja *JWTAuthenticator) BuildJWTClaims(user JWTUser, duration time.Duration) jwt.MapClaims {
	return jwt.MapClaims{
		"userId": user.ID,
		"email": user.Email,
		"exp": time.Now().Add(duration).Unix(),
		"iat": time.Now().Unix(),
	}
}

func (ja *JWTAuthenticator) GenerateToken(user JWTUser, exp time.Duration) (string, error) {
	claims := ja.BuildJWTClaims(user, exp)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenStr, err := token.SignedString([]byte(ja.secret))
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}

func (ja *JWTAuthenticator) ValidateToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method %v", t.Header["alg"])
		}
		return []byte(ja.secret), nil
	},
		jwt.WithExpirationRequired(),
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}),
	)
} 

func (ja *JWTAuthenticator) JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authHeader := r.Header.Get("Authorization")
		log.Println(
			strings.Contains(authHeader, "Bearer "),
		)
		if !strings.Contains(authHeader, "Bearer") {
			http.Error(w, "Invalid Token", http.StatusUnauthorized)
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := ja.ValidateToken(tokenStr)
		if err != nil {
			http.Error(w, "Bearer token malformed", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			http.Error(w, "Bearer token not contains user info", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(context.Background(), "user", claims)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}