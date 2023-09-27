package auth

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func FailOnError(err error, message string) {
	if err != nil {
		print("Token error: %v", message)
	}
}

type User struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

var TOKEN_SECRET = []byte("d5270e11-410c-452d-a4ca-66741df003a8")

func GenerateToken(user_id uuid.UUID, user User) string {
	authToken := jwt.MapClaims{
		"authorized": true,
		"user_id":    user_id,
		"user_name":  user.Username,
		"email":      user.Email,
		"exp":        time.Now().Add(time.Hour * 2),
	}

	AuthToken, _ := jwt.NewWithClaims(jwt.SigningMethodES256, authToken).SignedString(TOKEN_SECRET)
	return fmt.Sprintf("Bearer %v", AuthToken)
}

func TokenIsValid(cookie *http.Cookie) error {
	tokenString := GetTokenFromCookie(cookie)
	token, err := GetJWTToken(tokenString)

	FailOnError(err, "Error during validation")

	if err != nil {
		return err
	}

	if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return nil
	}

	return errors.New("token is not valid")
}

func GetTokenFromCookie(cookie *http.Cookie) string {
	token := cookie.Value

	if len(strings.Split(token, " ")) == 2 {
		return strings.Split(token, " ")[1]
	}
	return ""
}

func GetTokenIdFromCookie(cookie *http.Cookie) (uuid.UUID, error) {
	tokenString := GetTokenFromCookie(cookie)
	token, err := GetJWTToken(tokenString)

	FailOnError(err, "error getting token id")
	if err != nil {
		return uuid.Nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if ok && token.Valid {
		return uuid.MustParse(fmt.Sprintf("%v", claims["user_id"])), nil
	}

	return uuid.Nil, errors.New("error extracting ID")
}

func GetJWTToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("error during signin: %v", token.Header)
		}
		return TOKEN_SECRET, nil
	})

	return token, err
}

func HashPassword(password string) (string, error) {
	hashedValue, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedValue), err
}

func ValidatePassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
