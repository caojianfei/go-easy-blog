package helpers

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"go-easy-blog/libs"
	"net/http"
	"strings"
)

// 从token解析用户id
func GetUserIdFromAuthorization(tokenString string) (string, error) {
	secretKey := libs.GetConfig().JWTSecret
	token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (i interface{}, e error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(*jwt.StandardClaims)
	if !ok {
		return "", nil
	}
	if !token.Valid {
		return "", errors.New("token 过期")
	}

	return claims.Id, nil
}

// 从请求头中获取token
func GetAuthorizationTokenFormRequestHeader(request *http.Request) (string, error) {
	headers := request.Header
	Authorization, ok := headers["Authorization"]
	if !ok {
		return "", errors.New("authorization 不存在")
	}

	AuthorizationString := Authorization[0]
	if !strings.HasPrefix(AuthorizationString, "Bearer") {
		return "", errors.New("authorization 参数错误")
	}

	arr := strings.Split(AuthorizationString, " ")
	return arr[1], nil
}