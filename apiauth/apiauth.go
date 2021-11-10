package apiauth

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

var AUTHHEADER string = "X-Api-Key"

func Middleware(next http.Handler, db *sql.DB) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		token, err := GetTokenFromRequest(r)

		if err != nil {
			code := http.StatusNotFound
			http.Error(w, `{"message" : "token not found"}`, code)
			return
		}
		//Check if user exists
		ok := VerifyToken(token, db)
		if !ok {
			code := http.StatusUnauthorized
			http.Error(w, `{"message" : "Unauthorized"}`, code)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func GetTokenFromRequest(request *http.Request) (res string, err error) {
	result := ""
	result = request.Header.Get(AUTHHEADER)
	if result == "" {
		return result, errors.New("token not found")
	}
	return result, nil
}

func VerifyToken(tokenString string, db *sql.DB) bool {

	res := false

	fmt.Println(tokenString)
	fmt.Println(os.Getenv("WEB_TOKEN"))

	if tokenString == os.Getenv("WEB_TOKEN") {
		res = true
	}
	return res

}
