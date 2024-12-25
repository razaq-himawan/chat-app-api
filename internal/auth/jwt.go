package auth

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/razaq-himawan/chat-app-api/internal/app/model"
	"github.com/razaq-himawan/chat-app-api/utils"
)

type contextKey string

const UserKey contextKey = "userID"

var secret = []byte(os.Getenv("JWT_SECRET"))

func AuthJWT(userService model.UserService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenString := utils.GetTokenFromCookie(r)

			userID, err := GetUserIDFromToken(tokenString)
			if err != nil {
				permissionDenied(w)
				return
			}

			u, err := userService.GetUserByID(userID)
			if err != nil {
				log.Printf("failed to get user by id: %v", err)
				permissionDenied(w)
				return
			}

			ctx := context.WithValue(r.Context(), UserKey, u.ID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func CreateJWT(userID string) (string, error) {
	expiration := time.Second * time.Duration(3600*24*7)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":    userID,
		"expires_at": time.Now().Add(expiration).Unix(),
	})

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenString, err
}

func validateJWT(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secret), nil
	})
}

func permissionDenied(w http.ResponseWriter) {
	utils.WriteError(w, http.StatusForbidden, fmt.Errorf("permission denied"))
}

func GetUserIDFromContext(ctx context.Context) string {
	userID, ok := ctx.Value(UserKey).(string)
	if !ok {
		return ""
	}

	return userID
}

func GetUserIDFromToken(tokenString string) (string, error) {
	token, err := validateJWT(tokenString)
	if err != nil {
		return "", fmt.Errorf("failed to validate token: %v", err)
	}

	if !token.Valid {
		return "", fmt.Errorf("invalid token: %v", err)
	}

	claims := token.Claims.(jwt.MapClaims)
	userID := claims["user_id"].(string)

	return userID, nil
}
