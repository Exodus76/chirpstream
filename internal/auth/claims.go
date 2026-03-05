package auth

import (
	"context"
	"log"

	"github.com/golang-jwt/jwt/v5"
)

type CustomClaim struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}

func GetUserIdFromContext(ctx context.Context) (int, bool) {
	user, ok := ctx.Value("user").(*CustomClaim)
	if !ok {
		log.Printf("Error: failed getting user from context\n")
		return 0, false
	}
	return user.UserID, true
}
