package main

import (
	"fmt"
	"gosshc/internal/config"
	"gosshc/internal/services"
	"time"
)

func main() {

	privKey := []byte(config.Id_rsa)
	user := config.User
	ip := fmt.Sprintf("%s:%s", config.Ip, config.Port)

	sshC := services.NewSshClient(ip, user, privKey)
	for {

		err := sshC.SetupReverseTunnel("localhost:8080", "localhost:8080")
		if err != nil {
			fmt.Println(err)
			time.Sleep(5 * time.Second)
			continue
		}

		fmt.Println("tunnel up")

		err = <-sshC.ErrChan
		fmt.Println("💥:", err)

		sshC.Close()

	}

}
