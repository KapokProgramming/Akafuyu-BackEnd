package middleware

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"server/config"
	"server/model"
	"server/util"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func ValidateToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Add("Vary", "Authorization")

		authHeader := r.Header.Get("Authorization")

		headerPart := strings.Split(authHeader, " ")

		if len(headerPart) != 2 {
			// app.logger.Println("error :" + authHeader)
			w.WriteHeader(400)
			util.StandardResponseWriter(w, model.StandardResponse{
				Status: "fail",
				Data:   "invalid auth header",
			})
			return
		}

		if headerPart[0] != "Bearer" {
			w.WriteHeader(400)
			util.StandardResponseWriter(w, model.StandardResponse{
				Status: "fail",
				Data:   "invalid auth header",
			})
			return
		}

		token := headerPart[1]
		// do some check

		if len(token) <= 0 {
			w.WriteHeader(401)
			util.StandardResponseWriter(w, model.StandardResponse{
				Status: "fail",
				Data:   "unauthorized",
			})
			return
		}

		next.ServeHTTP(w, r)
	})
}

func ValidateJWT(signed_string string) (int, error) {
	token, err := jwt.Parse(signed_string, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	var user_id int
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		expiry, ok := claims["ExpiresAt"].(jwt.NumericDate)
		if !ok {
			return -1, fmt.Errorf("Unexpected ExpiresAt: %v", claims["ExpiresAt"])
		}
		notBefore, ok := claims["NotBefore"].(jwt.NumericDate)
		if !ok {
			return -1, fmt.Errorf("Unexpected NotBefore: %v", claims["NotBefore"])
		}
		if time.Now().Before(notBefore.Time) {
			return -1, fmt.Errorf("Invalid time: %v", notBefore)
		}
		if time.Now().After(expiry.Time) {
			return -1, fmt.Errorf("Expired: %v", expiry)
		}

		db, err := config.OpenDB(config.LoadConfig())
		if err != nil {
			return -1, fmt.Errorf("Failed to connect to db: %v", claims["Issuer"])
		}

		query := "SELECT user_id FROM users WHERE user_id=?;"
		err = db.QueryRow(query, claims["Issuer"]).Scan(&user_id)
		switch {
		case err == sql.ErrNoRows:
			return -1, fmt.Errorf("Invalid user_id: %v", claims["Issuer"])
		case err != nil:
			panic(err)
		}
		defer db.Close()
	} else {
		return -1, nil
	}
	if err != nil {
		return -1, err
	}
	if token.Valid {
		return -1, nil
	}
	return user_id, nil
}
