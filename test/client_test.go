package test

import (
	"fmt"
	"testing"
	"time"

	"github.com/ManuJSM/SshTunneler"
)

const (
	tunnelAddrLocal   = "localhost:8080"
	tunnelAddrReverse = "localhost:8081"
)

var tunnels = []SshTunneler.TunnelConfig{
	{LocalAddr: tunnelAddrLocal, RemoteAddr: tunnelAddrLocal, Reverse: false},
	{LocalAddr: tunnelAddrReverse, RemoteAddr: tunnelAddrReverse, Reverse: true},
}

const RECONNTIMEOUT = 3 * time.Second

func TestClient(t *testing.T) {

	Ip := "localhost"
	Port := "22"
	Id_rsa := ""
	User := "Pedro"

	privKey := []byte(Id_rsa)
	user := User
	ip := fmt.Sprintf("%s:%s", Ip, Port)

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
