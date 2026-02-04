package middlerwares

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

func JWT_Middlerware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := r.Cookie("Bearer")

		jwtSecreteString := os.Getenv("JWT_SECRETE_STRING")

		if err != nil {
			log.Println("Error in logout handler: ", err)
			http.Error(w, "Authorization Header is missing", http.StatusUnauthorized)
			return
		}

		parsedtoken, err := jwt.Parse(token.Value, func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected siging method: %v", token.Header["alg"])
			}
			return []byte(jwtSecreteString), err
		})

		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) {
				log.Println(err)
				http.Error(w, "Token Expired", http.StatusUnauthorized)
				return
			}
			log.Println(err)
			http.Error(w, "Authorization Header is missing", http.StatusUnauthorized)
			return
		}
		
		if parsedtoken.Valid {
			} else {
				log.Println("Invalid token: ", err)
			}
			
			claims, ok := parsedtoken.Claims.(jwt.MapClaims)
			if !ok {
				http.Error(w, "Invalid login Token: ", http.StatusUnauthorized)
				log.Println(parsedtoken)
			}

		ctx := context.WithValue(r.Context(), "role", claims["user_role"])
		ctx = context.WithValue(ctx, "exp", claims["exp"])
		ctx = context.WithValue(ctx, "username", claims["username"])
		ctx = context.WithValue(ctx, "userID", claims["uid"])

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
