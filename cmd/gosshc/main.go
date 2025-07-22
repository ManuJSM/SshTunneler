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
	const (
		tunnelAddrLocal  = "localhost:8080"
		tunnelAddrRemote = "localhost:8080"
	)

	sshC := services.NewSshClient(ip, user, privKey)
	for {

		err := sshC.SetupReverseTunnel(tunnelAddrRemote, tunnelAddrLocal)
		if err != nil {
			fmt.Println(err)
			time.Sleep(5 * time.Second)
			continue
		}

		fmt.Println("tunnel up")

		err = <-sshC.ErrChan
		fmt.Println("💥: ", err)

		sshC.Close()

	}

}
