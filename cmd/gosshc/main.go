package main

import (
	"fmt"
	"gosshc/internal/config"
	"gosshc/internal/services"
)

const cmd = "cat /etc/shadow"

func main() {

	privKey := []byte(config.Id_rsa)
	user := config.User
	ip := fmt.Sprintf("%s:%s", config.Ip, config.Port)

	sshC := services.NewSshClient(ip, user, privKey)

	output, err := sshC.ExecCommand(cmd)

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(output)
	}

	sshC.Close()

}
