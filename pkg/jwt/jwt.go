package jwt

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	. "github.com/cakra17/social/internal/utils"
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

type userClaimsKey struct{}

func NewJWTAuthenticator(secret string, duration time.Duration) *JWTAuthenticator {
	return &JWTAuthenticator{
		secret: secret,
		duration: duration,
	}
}

func (ja *JWTAuthenticator) BuildJWTClaims(user JWTUser) jwt.MapClaims {
	return jwt.MapClaims{
		"userId": user.ID,
		"email": user.Email,
		"exp": time.Now().Add(ja.duration).Unix(),
		"iat": time.Now().Unix(),
	}
}

func (ja *JWTAuthenticator) GenerateToken(user JWTUser) (string, error) {
	claims := ja.BuildJWTClaims(user)
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
		if !strings.Contains(authHeader, "Bearer") {
			JsonError(ErrNoTokenProvided, w)
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenStr == "" {
			JsonError(ErrTokenMalformed, w)
			return
		}

		token, err := ja.ValidateToken(tokenStr)
		if err != nil {
			JsonError(ErrTokenExpires, w)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			ne := CreateNewError(http.StatusUnauthorized, "Bearer token not contains user info")
			JsonError(ne, w)
			return
		}

		ctx := context.WithValue(r.Context(), userClaimsKey{}, claims)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (ja *JWTAuthenticator) GetClaims(ctx context.Context) (jwt.MapClaims, bool) {
	val := ctx.Value(userClaimsKey{})
	claims, ok := val.(jwt.MapClaims)
	return claims, ok
}