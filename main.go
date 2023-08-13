package main

import (
	"log"

	cli "movies/cli"
	
	

	dotenv "github.com/joho/godotenv"
)

func main() {

	if err := dotenv.Load("./.env"); err != nil {
		log.Fatal(err)
	}

	if err := cli.Execute(); err != nil {
		log.Fatal(err)
	}
}
