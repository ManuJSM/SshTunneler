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

func (c *Client) ConnectAndSetup(tunnels []TunnelConfig) error {
	if err := c.conn.Connect(); err != nil {
		return err
	}
	if err := c.tunnels.SetupTunnels(tunnels); err != nil {
		c.conn.Close()
		return fmt.Errorf("error seteando los tuneles %w", err)
	}
	c.monitor.Start()
	go c.watchErrors()
	return nil
}

func (c *Client) watchErrors() {
	err := <-c.monitor.OnError()
	log.Println("Monitor error:", err)

	for {
		if err = c.tryReconnect(); err == nil {
			break
		}
		log.Println(err)
		time.Sleep(5 * time.Second)
	}
}

func (c *Client) tryReconnect() error {
	c.Close()
	return c.ConnectAndSetup(c.tunnels.GetCurrentConf())
}

func (c *Client) Close() error {
	c.monitor.Stop()
	c.tunnels.CloseAll()
	return c.conn.Close()
}
