package main

import "github.com/chumaumenze/gjs"

type Data struct {
	Code    int
	Message string `json:"message"`
	inner   any
}

func main() {
	data := Data{
		Code:    200,
		Message: "Hello World!",
		inner:   "I am ignored",
	}
	_ = gjs.ValueOf(data) // e.g. {Code: 200, "message": "Hello World!"}
}
