package SshTunneler

import (
	"fmt"

	"github.com/ManuJSM/SshTunneler/config"
	"github.com/ManuJSM/SshTunneler/internal/services"
)

type Client struct {
	conn    services.SSHConnection
	tunnels services.TunnelManager
	monitor services.ConnectionMonitor
}

// NewClient crea un nuevo cliente SSH con la conexión, el monitor y el manejador de túneles configurados.
func NewClient(addr, user string, privKey []byte) *Client {
	conn := services.NewSSHConnection(addr, user, privKey)
	monitor := services.NewMonitor(conn)
	tunnels := services.NewTunnelManager(conn, monitor)
	return &Client{
		conn:    conn,
		tunnels: tunnels,
		monitor: monitor,
	}
}

// ConnectAndSetup establece la conexión SSH y configura los túneles proporcionados.
func (c *Client) ConnectAndSetup(tunnels ...config.TunnelConfig) error {
	if err := c.conn.Connect(); err != nil {
		return err
	}

	if len(tunnels) > 0 {
		if err := c.tunnels.SetupTunnels(tunnels); err != nil {
			c.conn.Close()
			return fmt.Errorf("error seteando los túneles %w", err)
		}
	}

	c.monitor.Start()

	return nil
}

// WatchErrors se queda esperando a que se generen errores en el monitor de conexión.
func (c *Client) WatchErrors() error {
	err := c.monitor.OnError()
	return fmt.Errorf("monitor error: %v", err)
}

func (c *Client) Close() error {
	c.monitor.Stop()
	c.tunnels.CloseAll()
	return c.conn.Close()
}
