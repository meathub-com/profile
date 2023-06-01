package transport

import (
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v4"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

// validateToken - validates an incoming JWT token
func validateToken(accessToken string) (map[string]interface{}, error) {
	var mySigningKey = []byte("missionimpossible")
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("could not validate auth token")
		}
		return mySigningKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("could not validate auth token")
}

// JWTAuth - a handy middleware function that will provide basic auth around specific endpoints
func JWTAuth(original func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			w.WriteHeader(http.StatusUnauthorized)
			log.Error("an unauthorized request has been made")
			return
		}

		authHeaderParts := strings.Split(authHeader, " ")
		if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != "bearer" {
			w.WriteHeader(http.StatusUnauthorized)
			log.Error("authorization header could not be parsed")
			return
		}

		claims, err := validateToken(authHeaderParts[1])
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			log.Error("could not validate incoming token")
			return
		}

		userID, ok := claims["sub"].(string)
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			log.Error("user ID not found in token claims")
			return
		}

		ctx := context.WithValue(r.Context(), "user_id", userID)
		original(w, r.WithContext(ctx))
	}
}
