package config

import (
	"log"

	"github.com/joho/godotenv"
)

func init() {

	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

}
