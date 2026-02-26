package helpers

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

func GetToken(ctx context.Context) (any, error) {
	ginCtx, _ := GinContextFromContext(ctx)
	authHeader := ginCtx.GetHeader("Authorization")

	if authHeader == "" {
		return nil, fmt.Errorf("no Authorization header provided")
	}

	// Pastikan authHeader adalah string
	header := authHeader

	// Memeriksa apakah header dimulai dengan "Bearer "
	if !strings.HasPrefix(header, "Bearer ") {
		return nil, fmt.Errorf("invalid Authorization header: must start with 'Bearer '")
	}

	// Mengambil token dengan menghapus "Bearer " dari header
	token := strings.TrimPrefix(header, "Bearer ")

	// Validasi token (opsional, tergantung kebutuhan)
	if token == "" {
		return nil, fmt.Errorf("token is empty")
	}

	// Gunakan token sesuai kebutuhan, misalnya untuk autentikasi atau mengambil data pengguna
	// Contoh: return user, nil
	return token, nil // Ganti dengan logika sesuai kebutuhan
}

func GetUserID(tokenStr string) (int, error) {
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
