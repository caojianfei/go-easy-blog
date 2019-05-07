package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/julienschmidt/httprouter"
	"github.com/siddontang/go-log/log"
	"go-easy-blog/database"
	"go-easy-blog/errors"
	"go-easy-blog/libs"
	"go-easy-blog/response"
	"net/http"
	"strconv"
	"time"
)

type UserCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Login(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	userCredentials := UserCredentials{}
	err := json.NewDecoder(r.Body).Decode(&userCredentials)
	if err != nil {
		log.Errorf("parse request body error, err: ", err.Error())
		response.SystemError(w)
		return
	}
	// 数据库检索用户
	db := database.New()
	stmt, err := db.Prepare(database.QueryUserByUsername)
	if err != nil {
		log.Errorf("prepare sql error, sql: %s, err: %s \n", database.QueryUserByUsername, err.Error())
		_ = response.SystemError(w)
		return
	}
	defer stmt.Close()

	user := &database.UserModel{}
	err = stmt.QueryRow(userCredentials.Username).Scan(&user.Id, &user.Username, &user.Password, &user.Nickname)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Info("login user is not exist")
			response.Fail(errors.UsernameOrPasswordError, "账号或密码错误", w)
			return
		} else {
			log.Errorf("query user error, sql format: %s, sql param: %s, err: %s", database.QueryUserByUsername, user.Username, err.Error())
			response.SystemError(w)
			return
		}
	}

	// 判断密码是否正确
	if userCredentials.Password != user.Password {
		response.Fail(errors.UsernameOrPasswordError, "账号或密码错误", w)
		return
	}

	// 生成登录token
	token, err := createJWTToken(user.Id)
	if err != nil {
		response.Fail(errors.SystemError, "登录失败", w)
		return
	}

	response.Success(map[string]string{"token":token, "nickname":user.Nickname}, w)
	return
}

type MyCustomClaims struct {
	UserId int `json:"userId"`
	jwt.StandardClaims
}


func createJWTToken(userId int) (string, error) {
	secret := libs.GetConfig().JWTSecret
	expire := libs.GetConfig().TokenExpire
	nowTime := time.Now()
	expire += nowTime.Unix()

	claims := &jwt.StandardClaims{Id:strconv.Itoa(userId), ExpiresAt: expire}
	fmt.Printf("login claims is %v \n", claims)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(secret))
	return ss, err
}
