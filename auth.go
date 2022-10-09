package main

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
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
		expiresAt := time.Unix(int64(claims["exp"].(float64)), 0)
		notBefore := time.Unix(int64(claims["nbf"].(float64)), 0)
		if time.Now().Before(notBefore) {
			return -1, fmt.Errorf("Invalid time: %v", notBefore)
		}
		if time.Now().After(expiresAt) {
			return -1, fmt.Errorf("Expired: %v", expiresAt)
		}
		db := createConnectionToDatabase()
		query := "SELECT user_id FROM users WHERE user_id=?;"
		err = db.QueryRow(query, claims["iss"]).Scan(&user_id)
		switch {
		case err == sql.ErrNoRows:
			return -1, fmt.Errorf("Invalid user_id: %v", claims["iss"])
		case err != nil:
			panic(err)
		}
	} else {
		return -1, fmt.Errorf("Invalid token: %v, ok: %v", token, ok)
	}
	if err != nil {
		return -1, err
	}
	return user_id, nil
}

func GetHashedPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func ValidatePassword(password string, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
