package crypt

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

//HashText generates hased password
func HashText(password string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes)
}

//CheckTextHash compares hash and user plain password
func CheckTextHash(text, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(text))
	if err != nil {
		return err == nil
	}
	return true
}

//Jwt generates json web tokens for stateless authentication
func Jwt(payload interface{}) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"_id":   payload,
		"exp":   time.Now().Add(time.Hour * 8760).Unix(),
		"admin": false,
	})
	tokenString, error := token.SignedString([]byte("secret"))
	if error != nil {
		fmt.Println(error)
	}
	return tokenString
}
