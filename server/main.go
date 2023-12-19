package main

import (
	"fmt"
	"os"
	"visio/internal/server"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	server := server.NewServer()
	fmt.Println("Server started on ", os.Getenv("PORT"))
	err = server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
