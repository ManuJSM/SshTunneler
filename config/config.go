package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	Id_rsa string
	Ip     string
	User   string
	Port   string
)

type TunnelConfig struct {
	LocalAddr  string
	RemoteAddr string
	Reverse    bool
}

func init() {

	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	Id_rsa = os.Getenv("id_rsa")
	Ip = os.Getenv("ip")
	User = os.Getenv("user")
	Port = os.Getenv("port")

}
