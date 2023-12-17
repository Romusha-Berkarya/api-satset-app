package middleware

import (
	"gateway-api-satset/entity"
	"gateway-api-satset/helper"
	"net/http"
	"strings"
)

const TOKEN_FAILED = "Token Invalid!"
const HEADER_ABSENT = "Authorization Header Absent!"
const ACCESS_DENIED = "Required Access!"
const WHITELIST_IP = "127.0.0.1"

func IsAuthorized(w http.ResponseWriter, r *http.Request) (bool, *entity.Token) {

	clientIp, _ := helper.GetIP(r)

	if clientIp != WHITELIST_IP {
		response := entity.ResponseMessage{
			Message: ACCESS_DENIED,
		}
		helper.Response(w, response, 401)

		return false, nil
	}

	Authorization := r.Header.Get("Authorization")
	bearerToken := strings.Split(Authorization, "Bearer ")

	if r.Header["Authorization"] != nil {
		token := validateToken(w, r, bearerToken[1])

		if token == nil {
			return false, nil
		}

		return true, token
	}

	return false, nil
}

func validateToken(w http.ResponseWriter, r *http.Request, bearer string) *entity.Token {
	var token entity.Token

	if err := helper.Instance.Table("tokens").Where("ID = ?", bearer).Take(&token).Error; err != nil {
		response := entity.ResponseMessage{
			Message: TOKEN_FAILED,
		}

		helper.Response(w, response, 401)

		return nil
	}

	return &token
}
