package services

import (
	"fmt"
	"io"
	"log"
	"net"
	"sync"

	"golang.org/x/crypto/ssh"
)

type TunnelManager interface {
	SetupTunnels(tunnels []TunnelConfig) error
	CloseAll() error
}

type TunnelConfig struct {
	LocalAddr  string
	RemoteAddr string
	Reverse    bool
}

type tunnelManagerImpl struct {
	conn       SSHConnection
	activeTuns []net.Listener
	monitor    ConnectionMonitor
}

func NewTunnelManager(conn SSHConnection, mon ConnectionMonitor) TunnelManager {
	return &tunnelManagerImpl{conn: conn, monitor: mon}
}

func (tunM *tunnelManagerImpl) SetupTunnels(tunnels []TunnelConfig) error {

	for _, t := range tunnels {
		var err error
		if t.Reverse {
			err = tunM.setupReverseTunnel(t.LocalAddr, t.RemoteAddr)
		} else {
			err = tunM.setupLocalTunnel(t.LocalAddr, t.RemoteAddr)
		}
		if err != nil {
			tunM.CloseAll()
			return err
		}
	}
	log.Println("Tunneles seteados")
	return nil
}
func (tunM *tunnelManagerImpl) CloseAll() (err error) {

	for _, t := range tunM.activeTuns {

		err = t.Close()

	}

	return err
}
func handleReverseTunnel(remoteConn net.Conn, localAddr string) {

	var wg sync.WaitGroup

	localConn, err := net.Dial(_SSHTYPECONN, localAddr)
	if err != nil {
		log.Printf("Failed to connect to local service: %v", err)
		remoteConn.Close()
		return
	}

	wg.Add(2)
	go func() {
		defer wg.Done()
		io.Copy(localConn, remoteConn)
	}()
	go func() {
		defer wg.Done()
		io.Copy(remoteConn, localConn)
	}()

	wg.Wait()
	localConn.Close()
	remoteConn.Close()

}

func (tunM *tunnelManagerImpl) setupReverseTunnel(remoteAddr, localAddr string) error {

	conn, err := tunM.conn.Client()

	if err != nil {
		return err
	}

	listener, err := conn.Listen(_SSHTYPECONN, remoteAddr)
	if err != nil {
		return fmt.Errorf("failed to set up remote listener: %w", err)
	}
	tunM.activeTuns = append(tunM.activeTuns, listener)

	go func() {

		defer func() {
			listener.Close()
			tunM.monitor.ReportError(err)
		}()

		for {
			var remoteConn net.Conn
			remoteConn, err = listener.Accept()
			if err != nil {
				return
			}
			go handleReverseTunnel(remoteConn, localAddr)
		}
	}()
	return nil
}

func handleLocalTunnel(remoteAddr string, localConn net.Conn, sshConn *ssh.Client) {

	var wg sync.WaitGroup

	remoteConn, err := sshConn.Dial(_SSHTYPECONN, remoteAddr)
	if err != nil {
		log.Printf("Failed to connect to remote service: %v", err)
		localConn.Close()
		return
	}

	wg.Add(2)
	go func() {
		defer wg.Done()
		io.Copy(localConn, remoteConn)
	}()
	go func() {
		defer wg.Done()
		io.Copy(remoteConn, localConn)
	}()

	wg.Wait()
	localConn.Close()
	remoteConn.Close()

}

func (tunnM *tunnelManagerImpl) setupLocalTunnel(remoteAddr, localAddr string) error {

	listener, err := net.Listen(_SSHTYPECONN, localAddr)
	if err != nil {
		return err
	}
	tunnM.activeTuns = append(tunnM.activeTuns, listener)
	go func() {
		defer func() {
			listener.Close()
			tunnM.monitor.ReportError(err)
		}()
		for {
			var localConn net.Conn
			var conn *ssh.Client

			localConn, err = listener.Accept()
			if err != nil {
				return
			}
			conn, err = tunnM.conn.Client()
			if err != nil {
				localConn.Close()
				return
			}
			go handleLocalTunnel(remoteAddr, localConn, conn)
		}
	}()

	return nil

}
