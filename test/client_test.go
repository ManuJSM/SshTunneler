package test

import (
	"fmt"
	"testing"
	"time"

	"github.com/ManuJSM/SshTunneler"
	"github.com/ManuJSM/SshTunneler/config"
)

const (
	tunnelAddrLocal   = "localhost:8080"
	tunnelAddrReverse = "localhost:8081"
)

var tunnels = []config.TunnelConfig{
	{LocalAddr: tunnelAddrLocal, RemoteAddr: tunnelAddrLocal, Reverse: false},
	{LocalAddr: tunnelAddrReverse, RemoteAddr: tunnelAddrReverse, Reverse: true},
}

const RECONNTIMEOUT = 3 * time.Second

func TestClient(t *testing.T) {

	privKey := []byte(config.Id_rsa)
	user := config.User
	ip := fmt.Sprintf("%s:%s", config.Ip, config.Port)

	sshC := SshTunneler.NewClient(ip, user, privKey)
	err := sshC.ConnectAndSetup(tunnels...)
	if err != nil {
		t.Log(err)
		return
	}
	err = sshC.WatchErrors()
	t.Log(err)
	sshC.Close()
}
