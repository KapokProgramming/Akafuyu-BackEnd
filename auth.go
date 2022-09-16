package main

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func CreateJWT(user_id int) (string, error) {
	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		NotBefore: jwt.NewNumericDate(time.Now()),
		Issuer:    strconv.Itoa(user_id),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}
	return ss, nil
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
		db := createConnectionToDatabase()
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
