package security

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/valyala/fasthttp"
)

const (
	JWT_KEY = "32ef9d54a97622d1175fd7a47661b76f"
)

func CreateToken(userID string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["sub"] = userID
	claims["iat"] = time.Now().Unix()
	claims["ver"] = 1
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	_token, err := token.SignedString([]byte(JWT_KEY))

	if err != nil {
		return "", err
	}

	return _token, nil
}

func ExtractUserID(r *fasthttp.Request) (uid uuid.UUID, err error) {
	token, _ := ExtractTokenObject(r)
	claims := token.Claims.(jwt.MapClaims)

	if claims["sub"] == nil {
		return uuid.Nil, errors.New("invalid token: user id not found")
	}

	userID := claims["sub"].(string)
	return uuid.Parse(userID)
}

func ExtractToken(r *fasthttp.Request) string {
	keys := r.URI().QueryArgs()
	token := string(keys.Peek("token"))

	if token != "" {
		return token
	}

	bearerToken := string(r.Header.Peek("Authorization"))
	if len(strings.Split(bearerToken, " ")) == 2 {
		return strings.Split(bearerToken, " ")[1]
	}
	return ""
}

func ExtractTokenObject(r *fasthttp.Request) (*jwt.Token, error) {
	tokenString := ExtractToken(r)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(JWT_KEY), nil
	})

	if err != nil {
		return nil, err
	}
	return token, err
}

func Pretty(data interface{}) {
	b, err := json.MarshalIndent(data, "", " ")

	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println(string(b))
}

func GetSecret() string {
	return JWT_KEY
}
