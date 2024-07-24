package helper

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/valyala/fasthttp"
)

type TokenType struct {
	name   string
	exp    int
	secret string
}

type AuthInterface interface {
	CreateToken(userID string, tokenType string) (string, error)
	ExtractUserID(token string, tokenType string) (uuid.UUID, error)
	ExtractBearerToken(r *fasthttp.Request) string
}

type auth struct{}

func NewAuth() AuthInterface {
	return &auth{}
}

func (a *auth) CheckTokenType(tokenType string) TokenType {
	accessTokenType := TokenType{"access", 1, os.Getenv("JWT_ACCESS_SECRET")}

	switch tokenType {
	case "access":
		return accessTokenType
	case "refresh":
		return TokenType{"refresh", 168, os.Getenv("JWT_REFRESH_SECRET")}
	default:
		return accessTokenType
	}
}

func (a *auth) CreateToken(userId string, tokenType string) (string, error) {
	tType := a.CheckTokenType(tokenType)

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["sub"] = userId
	claims["iat"] = time.Now().Unix()
	claims["ver"] = 1
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(tType.exp)).Unix()

	_token, err := token.SignedString([]byte(tType.secret))

	if err != nil {
		return "", err
	}

	return _token, nil
}

func (a *auth) ExtractUserID(token string, tokenType string) (uid uuid.UUID, err error) {
	tType := a.CheckTokenType(tokenType)
	tokenObj, err := a.ExtractTokenObject(token, tType.secret)

	if err != nil {
		return uuid.Nil, err
	}

	claims := tokenObj.Claims.(jwt.MapClaims)

	if claims["sub"] == nil {
		return uuid.Nil, errors.New("invalid token: user id not found")
	}

	userID := claims["sub"].(string)
	return uuid.Parse(userID)
}

func (a *auth) ExtractBearerToken(r *fasthttp.Request) string {
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

func (a *auth) ExtractTokenObject(tokenString string, secret string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}
	return token, err
}

func (a *auth) Pretty(data interface{}) {
	b, err := json.MarshalIndent(data, "", " ")

	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println(string(b))
}
