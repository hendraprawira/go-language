package main

import (
	"fmt"
	"log"
	"os"

	"Remember-Golang/app/router"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println(err)
	}
	port := ":" + os.Getenv("ACTIVE_PORT")
	if err := router.Routes().Run(port); err != nil {
		log.Fatalln(err)
	}
}
