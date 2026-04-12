package utils

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPassword(hashedPassword, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)) == nil
}

type AuthClaims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func GenerateJWT(userID uuid.UUID, role string, secret []byte) (string, error) {
	claims := AuthClaims{
		UserID: userID.String(),
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID.String(),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

func ParseJWT(tokenStr string, secret []byte) (*AuthClaims, error) {
	claims := &AuthClaims{}

	_, err := jwt.ParseWithClaims(
		tokenStr,
		claims,
		func(token *jwt.Token) (any, error) {
			return secret, nil
		},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
	)
	if err != nil {
		return nil, err
	}

	if claims.UserID == "" || claims.Role == "" {
		return nil, fmt.Errorf("invalid token payload")
	}

	return claims, nil
}

func JWTAuthMiddleware(jwtSecret string) gin.HandlerFunc {
	secret := []byte(jwtSecret)

	return func(c *gin.Context) {
		tokenStr := strings.TrimSpace(c.GetHeader("Authorization"))
		if tokenStr == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token required"})
			return
		}

		tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")
		if tokenStr == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header"})
			return
		}

		claims, err := ParseJWT(tokenStr, secret)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		c.Set("userID", claims.UserID)
		c.Set("role", claims.Role)
		c.Next()
	}
}

func RoleMiddleware(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleRaw, exists := c.Get("role")
		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "role not found in token"})
			return
		}

		role, ok := roleRaw.(string)
		if !ok || role != requiredRole {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "insufficient role"})
			return
		}

		c.Next()
	}
}

type rateLimiterEntry struct {
	Count     int
	ResetTime time.Time
}

type RateLimiter struct {
	mu       sync.Mutex
	max      int
	window   time.Duration
	visitors map[string]*rateLimiterEntry
	secret   []byte
}

func NewRateLimiter(jwtSecret string, max int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		max:      max,
		window:   window,
		visitors: make(map[string]*rateLimiterEntry),
		secret:   []byte(jwtSecret),
	}
}

func (r *RateLimiter) key(c *gin.Context) string {
	if uidRaw, ok := c.Get("userID"); ok {
		if uid, ok := uidRaw.(string); ok && uid != "" {
			return "user:" + uid
		}
	}

	authHeader := strings.TrimSpace(c.GetHeader("Authorization"))
	authHeader = strings.TrimPrefix(authHeader, "Bearer ")
	if authHeader != "" {
		if claims, err := ParseJWT(authHeader, r.secret); err == nil && claims.UserID != "" {
			return "user:" + claims.UserID
		}
	}

	return "ip:" + c.ClientIP()
}

func (r *RateLimiter) Allow(key string) (bool, int, time.Time) {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	entry, exists := r.visitors[key]
	if !exists || now.After(entry.ResetTime) {
		entry = &rateLimiterEntry{Count: 0, ResetTime: now.Add(r.window)}
		r.visitors[key] = entry
	}

	entry.Count++
	remaining := r.max - entry.Count
	if remaining < 0 {
		remaining = 0
	}

	if entry.Count > r.max {
		return false, 0, entry.ResetTime
	}

	return true, remaining, entry.ResetTime
}

func RateLimiterMiddleware(jwtSecret string, max int, window time.Duration) gin.HandlerFunc {
	limiter := NewRateLimiter(jwtSecret, max, window)

	return func(c *gin.Context) {
		key := limiter.key(c)
		allowed, remaining, reset := limiter.Allow(key)

		c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", max))
		c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", remaining))
		c.Header("X-RateLimit-Reset", fmt.Sprintf("%d", reset.Unix()))

		if !allowed {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "rate limit exceeded"})
			return
		}

		c.Next()
	}
}
