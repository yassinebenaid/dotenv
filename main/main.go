package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	e, err := godotenv.Read()
	if err != nil {
		panic(err)
	}
	fmt.Printf("a = \"%s\"\n", e["a"])
	os.Expand(`hello\$TEST`, func(s string) string {
		fmt.Println(s)

		return ""
	})
}
