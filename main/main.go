package main

import (
	"fmt"

	"github.com/joho/godotenv"
)

func main() {

	e, err := godotenv.Read()
	if err != nil {
		panic(err)
	}
	fmt.Println(e)
}
