# SSH Tunnel Client

A robust SSH client written in Go for creating and managing SSH tunnels with monitoring capabilities.

## Features

- ✅ SSH connection with private key authentication
- ✅ Multiple SSH tunnel management
- ✅ Real-time connection monitoring
- ✅ DotEnv Configuration

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

	privKey := []byte(config.Id_rsa)
	user := config.User
	ip := fmt.Sprintf("%s:%s", config.Ip, config.Port)

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
