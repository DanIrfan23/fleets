package middlewares

import (
	"context"
	"fleets/models"
	"fleets/responses"
	"log"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt"
)

func JwtMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("jwt_token")
		if err != nil {
			log.Println(err)
			responses.ErrorResponse(w, http.StatusUnauthorized, "Unauthorized")

			return
		}

		token, err := jwt.ParseWithClaims(cookie.Value, &models.Claims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET_KEY")), nil
		})
		if err != nil || !token.Valid {
			log.Println(err)
			responses.ErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		claims, ok := token.Claims.(*models.Claims)
		if !ok {
			log.Println("Error extracting claims")
			responses.ErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		ctx := context.WithValue(r.Context(), "username", claims.Username)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
