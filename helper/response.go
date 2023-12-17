package helper

import (
	"encoding/json"
	"net/http"
)

func ReadFromRequestBody(request *http.Request, result interface{}) {
	decoder := json.NewDecoder(request.Body)
	decoder.Decode(result)
}

func Response(writer http.ResponseWriter, response interface{}, statusCode int) {
	writer.Header().Set("Content-Type", "application/json")

	writer.WriteHeader(statusCode)

	if response != nil {

		err := json.NewEncoder(writer).Encode(response)

		PanicIfError(err)
	}
}
