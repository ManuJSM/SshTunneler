package main

import (
	"fmt"
	"gosshc/internal/config"
	"gosshc/internal/services"
	"log"
	"time"
)

const (
	tunnelAddrLocal   = "localhost:8080"
	tunnelAddrReverse = "localhost:8081"
)

var tunnels = []services.TunnelConfig{
	{LocalAddr: tunnelAddrLocal, RemoteAddr: tunnelAddrLocal, Reverse: false},
	{LocalAddr: tunnelAddrReverse, RemoteAddr: tunnelAddrReverse, Reverse: true},
}

func main() {

	privKey := []byte(config.Id_rsa)
	user := config.User
	ip := fmt.Sprintf("%s:%s", config.Ip, config.Port)

	for {
		sshC := services.NewClient(ip, user, privKey)
		for {
			err := sshC.ConnectAndSetup(tunnels...)
			if err != nil {
				log.Println(err)
				time.Sleep(services.RECONNTIMEOUT)
			} else {
				break
			}
		}
		sshC.WatchErrors()
		sshC.Close()
	}

}
