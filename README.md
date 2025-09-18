# SSH Tunnel Client

A robust SSH client written in Go for creating and managing SSH tunnels with monitoring capabilities.

## Features

- ✅ SSH connection with private key authentication
- ✅ Multiple SSH tunnel management
- ✅ Real-time connection monitoring

---

## Instalation

```bash
go get github.com/ManuJSM/SshTunneler
```

## ⚙️ Basic Usage

```go
const (
	tunnelAddrLocal   = "localhost:8080"
	tunnelAddrReverse = "localhost:8081"
)

var tunnels = []config.TunnelConfig{
	{LocalAddr: tunnelAddrLocal, RemoteAddr: tunnelAddrLocal, Reverse: false},
	{LocalAddr: tunnelAddrReverse, RemoteAddr: tunnelAddrReverse, Reverse: true},
}

func main() {
	Ip := "localhost"
	Port := "22"
	Id_rsa := ""
	user := "Pedro"

	privKey := []byte(Id_rsa)
	ip := fmt.Sprintf("%s:%s", Ip, Port)

	client := NewClient(ip, user, privKey)
  err := client.ConnectAndSetup(tunnels...)
  if err != nil {
    panic(err)
	}

	err := client.WatchErrors()
  Log(err)
	sshC.Close()
}

```
