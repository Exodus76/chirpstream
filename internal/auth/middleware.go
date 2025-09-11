package auth

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/julienschmidt/httprouter"
)

func AuthMiddleware(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

		authHeader := strings.Split(r.Header.Get("Authorization"), "Bearer ")
		if len(authHeader) != 2 {
			log.Printf("Malformed token")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Malformed Token"))
		}

		jwtToken := authHeader[1]
		token, err := jwt.ParseWithClaims(jwtToken, &CustomClaim{}, func(t *jwt.Token) (any, error) {
			//TODO: remove the secret key from here (IMP)
			return []byte("mykey"), nil
		})

		if claims, ok := token.Claims.(*CustomClaim); ok && token.Valid {
			ctx := context.WithValue(r.Context(), "user", claims)
			next(w, r.WithContext(ctx), p)
		} else {
			log.Printf("ERROR: Unauthorized access %v\n", err)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized"))
		}

	}
}
