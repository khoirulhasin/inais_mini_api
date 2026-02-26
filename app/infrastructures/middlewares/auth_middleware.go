package middlewares

import (
	"context"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/khoirulhasin/untirta_api/app/models"

	"github.com/vektah/gqlparser/v2/gqlerror"
	"gorm.io/gorm"
)

var userCtxKey = &contextKey{"user"}

type contextKey struct {
	name string
}

// User represents a database-backed user object
type User struct {
	ID    int
	UUID  string
	Name  string
	Email string
	Roles []string // For role-based validation
}

// Middleware validates the Authorization header and packs the user into context
func AuthMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			log.Printf("No Authorization header found, allowing unauthenticated access")
			c.Next()
			return
		}

		// Check for Bearer token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.Next()
			return
		}

		tokenStr := parts[1]
		if tokenStr == "" {
			log.Printf("Empty Bearer token")
			c.Next()
			return
		}

		// Validate JWT and get user ID
		userID, err := validateAndGetUserID(tokenStr)
		if err != nil {
			log.Printf("Invalid token: %v", err)
			c.Next()
			return
		}

		// Get user from database
		user, err := getUserByID(db, userID)
		if err != nil {
			log.Printf("Failed to fetch user %d: %v", userID, err)
			c.Next()
			return
		}

		if user == nil {
			log.Printf("User %d not found in database", userID)
			c.Next()
			return
		}

		// Put user in context
		ctx := context.WithValue(c.Request.Context(), userCtxKey, user)
		c.Request = c.Request.WithContext(ctx)
		log.Printf("Authenticated user %d (%s) with roles: %v", user.ID, user.Name, user.Roles)

		// Proceed to next handler
		c.Next()
	}
}

// validateAndGetUserID validates the JWT token and extracts the user ID
func validateAndGetUserID(tokenStr string) (int, error) {
	// Replace with your JWT secret
	secretKey := []byte(os.Getenv("SECRET_KEY"))

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, gqlerror.Errorf("Metode penandatanganan tidak valid: %v", token.Header["alg"])
		}
		return secretKey, nil
	})

	if err != nil {
		return 0, err
	}

	if !token.Valid {
		return 0, gqlerror.Errorf("Token JWT tidak valid")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, gqlerror.Errorf("Klaim JWT tidak valid")
	}

	userIDFloat, ok := claims["ID"].(float64) // JSON numbers are float64
	if !ok {
		return 0, gqlerror.Errorf("user_id tidak valid di JWT")
	}

	userID := int(userIDFloat)
	if userID <= 0 {
		return 0, gqlerror.Errorf("user_id tidak valid: %d", userID)
	}

	// Check token expiration
	if exp, ok := claims["exp"].(float64); ok {
		if float64(time.Now().Unix()) > exp {
			return 0, gqlerror.Errorf("Token telah kedaluwarsa")
		}
	}

	return userID, nil
}

// getUserByID fetches a user and their roles from the database
func getUserByID(db *gorm.DB, userID int) (*User, error) {
	var user models.User
	err := db.Where("id = ?", userID).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // User not found
		}
		return nil, gqlerror.Errorf(err.Error())
	}

	// Fetch roles
	users2roles := []*models.Users2role{}
	err = db.Where("user_id = ?", userID).Preload("Role").Find(&users2roles).Error
	if err != nil {
		return nil, gqlerror.Errorf(err.Error())
	}

	roles := make([]string, 0, len(users2roles))
	for _, users2role := range users2roles {
		if users2role.Role != nil {
			roles = append(roles, users2role.Role.Code)
		} else {
			log.Printf("Warning: Role is nil for users2role with user_id=%d", userID)
		}
	}

	response := &User{
		ID:    user.ID,
		UUID:  user.UUID.String(),
		Name:  user.Username, // Map Username to Name
		Email: user.Email,
		Roles: roles,
	}

	return response, nil
}

// ForContext finds the user from the context
func ForContext(ctx context.Context) *User {
	raw, ok := ctx.Value(userCtxKey).(*User)
	if !ok {
		log.Printf("No user found in context or invalid type")
		return nil
	}
	return raw
}
