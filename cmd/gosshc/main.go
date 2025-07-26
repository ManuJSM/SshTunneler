package main

import (
	"fmt"
	"gosshc/internal/config"
	"gosshc/internal/services"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const (
	tunnelAddrLocal   = "localhost:8080"
	tunnelAddrReverse = "localhost:8081"
)

var tunnels = []services.TunnelConfig{
	{LocalAddr: tunnelAddrLocal, RemoteAddr: tunnelAddrLocal, Reverse: false},
	{LocalAddr: tunnelAddrReverse, RemoteAddr: tunnelAddrReverse, Reverse: true},
}

func trap_c(cleanup func()) {
	// Canal para recibir señales del sistema
	sigCh := make(chan os.Signal, 1)

	// Registramos las señales que queremos capturar
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	// Bloqueamos aquí hasta que llegue una señal
	sig := <-sigCh
	log.Printf("Señal recibida: %s. Iniciando limpieza...\n", sig)

	// Ejecutamos la función de limpieza que pasaron
	if cleanup != nil {
		cleanup()
	}

	log.Println("Limpieza completa. Saliendo.")
}

func main() {

	privKey := []byte(config.Id_rsa)
	user := config.User
	ip := fmt.Sprintf("%s:%s", config.Ip, config.Port)

	sshC := services.NewClient(ip, user, privKey)

	err := sshC.ConnectAndSetup(tunnels)
	if err != nil {
		log.Println(err)
		sshC.Close()
	}

	trap_c(func() { sshC.Close() })

}
