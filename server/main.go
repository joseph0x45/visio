package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"visio/internal/server"
)

func main() {
	appEnv := os.Getenv("ENV")
	if appEnv != "PROD" {
		err := godotenv.Load()
		if err != nil {
			panic(err)
		}
	}
	server := server.NewServer()
	fmt.Println("Server started on ", os.Getenv("PORT"))
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
