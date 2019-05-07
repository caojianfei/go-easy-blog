package middlewares

import (
	"database/sql"
	"github.com/dgrijalva/jwt-go"
	"github.com/julienschmidt/httprouter"
	"go-easy-blog/database"
	"go-easy-blog/libs"
	"io"
	"log"
	"net/http"
	"strings"
)

func CheckLogin(next httprouter.Handle) httprouter.Handle {

	return func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		headers := request.Header
		Authorization, ok := headers["Authorization"]
		if !ok {
			_, err :=notLoginError(writer)
			if err != nil {
				log.Fatal(err)
			}
			return
		}

		AuthorizationString := Authorization[0]
		if !strings.HasPrefix(AuthorizationString, "Bearer") {
			_, err := notLoginError(writer)
			if err != nil {
				log.Fatal(err)
			}
			return
		}

		arr := strings.Split(AuthorizationString, " ")
		tokenString := arr[1]
		// log.Println("tokenString:", tokenString)

		secretKey := libs.GetConfig().JWTSecret
		// fmt.Println("secret key is", secretKey)
		token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (i interface{}, e error) {
			return []byte(secretKey), nil
		})

		if err != nil {
			// log.Fatalf("error, err: %s", err)
			_, err := notLoginError(writer)
			if err != nil {
				log.Fatal(err)
			}
			return
		}

		claims, ok := token.Claims.(*jwt.StandardClaims)
		if !ok {
			//log.Fatalf("parse token fail")
			_, err := notLoginError(writer)
			if err != nil {
				log.Fatal(err)
			}
			return
		}
		if !token.Valid {
			// log.Fatalf("token is invalid")
			_, err := notLoginError(writer)
			if err != nil {
				log.Fatal(err)
			}
			return
		}

		// fmt.Printf("claims is %v \n", claims)

		userId := claims.Id

		// 根据用户id从数据库查询用户
		db := database.New()
		stmt, err := db.Prepare(database.QueryUserByUserId)
		if err != nil {
			// log.Fatalf("prepare sql fail, sql: %s, err : %s\n", database.QueryUserByUserId, err)
			_, err := notLoginError(writer)
			if err != nil {
				log.Fatal(err)
			}
			return
		}

		defer stmt.Close()
		user := database.UserModel{}
		err = stmt.QueryRow(userId).Scan(&user.Id, &user.Username, &user.Nickname)
		if err != nil {
			if err == sql.ErrNoRows {
				_, err := notLoginError(writer)
				if err != nil {
					log.Fatal(err)
				}
				return
			} else {
				_, err := notLoginError(writer)
				if err != nil {
					log.Fatal(err)
				}
				return
			}
		}
		// fmt.Printf("login user is %v \n", user)
		next(writer, request, params)
	}
}

func notLoginError(w http.ResponseWriter)(int, error) {
	w.WriteHeader(http.StatusUnauthorized)
	n, err := io.WriteString(w, "Unauthorized")
	return n, err
}
