package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	auth "shroom-wiki-backend/Auth"
)

func JSON(writter http.ResponseWriter, statusCode int, data interface{}) {
	writter.WriteHeader(statusCode)

	if err := json.NewEncoder(writter).Encode(data); err != nil {
		fmt.Fprintf(writter, "%s", err.Error())
	}
}

func Authentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writter http.ResponseWriter, request *http.Request) {
		token, err := request.Cookie("jwt")
		if err != nil || auth.TokenIsValid(token) != nil {
			JSON(writter, http.StatusUnauthorized, struct {
				Error string `json:"error"`
			}{
				Error: "unauthorized",
			})

			return
		}

		next.ServeHTTP(writter, request)
	})
}
