package helper

import (
	"fmt"
)

func PanicIfError(err error) {
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}

func ErrorWithMessage(err error, msg string) {
	if err != nil {
		fmt.Println(err)
		fmt.Println("message: ", msg)
		panic(err)
	}
}

func ErrorIfDataEmpty(data interface{}) bool {
	if data == nil {
		return true
	}

	return false
}
