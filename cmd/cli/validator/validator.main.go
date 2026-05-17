package main

import (
	"fmt"
	"log"

	"github.com/go-playground/validator/v10"
)

type User struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"` // xxx@xxx.com
	Age   int    `json:"age" validate:"gte=18"`
}

func main() {
	u := User{
		// Name:  "Phuc",
		Email: "phuc@phuc",
		Age:   17,
	}

	validate := validator.New()
	err := validate.Struct(u)
	if err != nil {
		fmt.Println("Validation failed!")

		for _, e := range err.(validator.ValidationErrors) {
			fmt.Printf("Field %s, Error %s\n", e.Field(), e.Tag())
		}
	} else {
		log.Println("Validation success!")
	}
}
