package handler

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"rip/internal/app/ds"
	"rip/internal/app/redis"
)

type contextKey string

const (
	UserIDKey contextKey = "user_id"
	RoleKey   contextKey = "role"
	LoginKey  contextKey = "login"
	jwtPrefix            = "Bearer "
)

func getJWTSecret() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "super-secret-jwt-key-historypigments-2026"
	}
	return []byte(secret)
}

// Генерация JWT токена
func GenerateJWTToken(userID uint, login, role string) (string, time.Duration, error) {
	expiresIn := 24 * time.Hour
	claims := ds.JWTClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "historypigments-api",
		},
		UserID: userID,
		Login:  login,
		Role:   role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(getJWTSecret())
	return tokenStr, expiresIn, err
}

// parseBearerToken: извлекает и проверяет JWT из заголовка Authorization (с учётом Redis Blacklist)
func parseBearerToken(r *http.Request, redisClient *redis.Client) (*ds.JWTClaims, string, int, string) {
	authHeader := strings.TrimSpace(r.Header.Get("Authorization"))
	if authHeader == "" {
		return nil, "", http.StatusUnauthorized, "требуется авторизация (заголовок 'Authorization: Bearer <token>')"
	}

	tokenStr := authHeader
	if strings.HasPrefix(authHeader, "Bearer ") {
		tokenStr = strings.TrimPrefix(authHeader, "Bearer ")
	} else if strings.HasPrefix(authHeader, "bearer ") {
		tokenStr = strings.TrimPrefix(authHeader, "bearer ")
	}
	tokenStr = strings.TrimSpace(tokenStr)

	if tokenStr == "" {
		return nil, "", http.StatusUnauthorized, "требуется авторизация (пустой токен)"
	}

	// Проверка токена в Blacklist Redis
	if redisClient != nil && redisClient.CheckJWTInBlacklist(r.Context(), tokenStr) {
		return nil, "", http.StatusUnauthorized, "токен отозван (выполнен logout)"
	}

	claims := &ds.JWTClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return getJWTSecret(), nil
	})
	if err != nil || !token.Valid {
		return nil, "", http.StatusUnauthorized, "недействительный или истекший JWT токен"
	}
	return claims, tokenStr, 0, ""
}

func withClaims(r *http.Request, claims *ds.JWTClaims) *http.Request {
	ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
	ctx = context.WithValue(ctx, RoleKey, claims.Role)
	ctx = context.WithValue(ctx, LoginKey, claims.Login)
	return r.WithContext(ctx)
}

// Auth Middleware: Проверка JWT токена и Redis Blacklist (доступ только создателю)
func (h *Handler) AuthMiddleware(redisClient *redis.Client, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, _, status, msg := parseBearerToken(r, redisClient)
		if claims == nil {
			writeErrorJSON(w, status, msg)
			return
		}
		next(w, withClaims(r, claims))
	}
}

// OptionalAuthMiddleware: гость проходит без токена, а с валидным токеном — как создатель
// (нужно для признака is_creator в списке и is_like в ленте)
func (h *Handler) OptionalAuthMiddleware(redisClient *redis.Client, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if claims, _, _, _ := parseBearerToken(r, redisClient); claims != nil {
			r = withClaims(r, claims)
		}
		next(w, r)
	}
}

// Извлечение текущего UserID из контекста запроса (0 — гость)
func GetUserIDFromContext(r *http.Request) uint {
	if val := r.Context().Value(UserIDKey); val != nil {
		if id, ok := val.(uint); ok && id > 0 {
			return id
		}
	}
	return 0
}

func GetUserRoleFromContext(r *http.Request) string {
	if val := r.Context().Value(RoleKey); val != nil {
		if role, ok := val.(string); ok {
			return role
		}
	}
	return "guest"
}
