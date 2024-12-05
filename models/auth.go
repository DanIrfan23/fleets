package models

import (
	"database/sql"
	"errors"
	"fleets/configs"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

type LoginDTO struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type TokenReturn struct {
	ExpirationTime time.Time
	Token          string
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func CheckAccountQuery(data *LoginDTO) (string, error) {
	db := configs.GetDB()
	var username string
	var password string

	query := "SELECT usrid, usrpass FROM cltbuser WHERE usrid = ?"
	err := db.QueryRow(query, data.Username).Scan(&username, &password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", errors.New("Username belum terdaftar")
		}
		return "", err
	}

	if password != data.Password {
		return "", errors.New("Password salah")
	}

	return username, nil
}

func GenerateTokenQuery(username string) (TokenReturn, error) {
	expirationTime := time.Now().Add(60 * time.Minute)
	claims := &Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		log.Println(err)
		return TokenReturn{}, fmt.Errorf("Failed to generate token")
	}

	tokenReturn := TokenReturn{
		Token:          tokenString,
		ExpirationTime: expirationTime,
	}

	return tokenReturn, err
}
