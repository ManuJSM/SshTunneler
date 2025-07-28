package services

import (
	"fmt"
	"log"
	"time"
)

type Client struct {
	conn    SSHConnection
	tunnels TunnelManager
	monitor ConnectionMonitor
}

const RECONNTIMEOUT = 30 * time.Second

func NewClient(addr, user string, privKey []byte) *Client {
	conn := NewSSHConnection(addr, user, privKey)
	monitor := NewMonitor(conn)
	tunnels := NewTunnelManager(conn, monitor)
	return &Client{
		conn:    conn,
		tunnels: tunnels,
		monitor: monitor,
	}
}

func (c *Client) ConnectAndSetup(tunnels ...TunnelConfig) error {
	if err := c.conn.Connect(); err != nil {
		return err
	}

	if len(tunnels) > 0 {
		if err := c.tunnels.SetupTunnels(tunnels); err != nil {
			c.conn.Close()
			return fmt.Errorf("error seteando los tuneles %w", err)
		}
	}

	c.monitor.Start()

	return nil
}

//se queda a la espera de errores cuando el monitor esta encendido

func (c *Client) WatchErrors() {
	err := c.monitor.OnError()
	log.Println("Monitor error:", err)
}

func (c *Client) Close() error {
	c.monitor.Stop()
	c.tunnels.CloseAll()
	return c.conn.Close()
}
