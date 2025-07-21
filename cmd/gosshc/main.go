package main

import (
	"fmt"
	"gosshc/internal/config"
	"gosshc/internal/services"
	"os"
	"os/signal"
	"syscall"
)

func waitForInterrupt() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c
}

func main() {

	privKey := []byte(config.Id_rsa)
	user := config.User
	ip := fmt.Sprintf("%s:%s", config.Ip, config.Port)

	sshC := services.NewSshClient(ip, user, privKey)

	err := sshC.SetupReverseTunnel("localhost:8080", "localhost:8080")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Túnel activo. Presiona Ctrl+C para salir.")
	waitForInterrupt()

	sshC.Close()

}
