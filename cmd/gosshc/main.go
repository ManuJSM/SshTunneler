package main

import (
	"fmt"
	"gosshc/internal/config"
	"gosshc/internal/services"
	"log"
	"time"
)

func main() {

	privKey := []byte(config.Id_rsa)
	user := config.User
	ip := fmt.Sprintf("%s:%s", config.Ip, config.Port)
	const _KEEPALIVETIMEOUT = 15 * time.Second

	const (
		tunnelAddrLocal   = "localhost:8080"
		tunnelAddrReverse = "localhost:8081"
	)
	sshC := services.NewSshClient(ip, user, privKey)
	ticker := time.NewTicker(_KEEPALIVETIMEOUT)
	for {

		err := sshC.SetupReverseTunnel(tunnelAddrReverse, tunnelAddrReverse)
		if err != nil {
			log.Println(err)
			time.Sleep(_KEEPALIVETIMEOUT)
		}

		err = sshC.SetupLocalTunnel(tunnelAddrLocal, tunnelAddrLocal)
		if err != nil {
			log.Println(err)
			time.Sleep(_KEEPALIVETIMEOUT)
		}

		if err == nil {

			log.Println("tunnels up")

			connection := true

			for connection {
				select {
				case err := <-sshC.ErrChan:
					log.Println(err)
					connection = false
				case <-ticker.C:
					connection = sshC.TestConnection()
					if !connection {
						log.Println("conexion caida ):")
					}
				}
			}

		}

		sshC.Close()
	}

}
