package services

import (
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/ssh"
)

const _SSHTYPECONN = "tcp"
const _TCPTIMEOUT = 1 * time.Second
const _KEEPALIVETIMEOUT = 5 * time.Second

type SshClient struct {
	addr    string
	user    string
	privKey []byte
	conn    *ssh.Client
	ErrChan chan error
}

func (sshC *SshClient) getAuthMethod() (ssh.AuthMethod, error) {

	signer, err := ssh.ParsePrivateKey(sshC.privKey)
	if err != nil {
		return nil, err
	}

	return ssh.PublicKeys(signer), nil

}

func (sshC *SshClient) getConn() (conn *ssh.Client, err error) {

	conn = sshC.conn

	// se comprueba si ya esta establecida la conexion
	if conn == nil {

		var authMethod ssh.AuthMethod

		authMethod, err = sshC.getAuthMethod()
		if err != nil {
			return
		}

		conn, err = ssh.Dial(_SSHTYPECONN, sshC.addr, &ssh.ClientConfig{
			User:            sshC.user,
			Auth:            []ssh.AuthMethod{authMethod},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(), // ⚠️ Conveniente cambiarlo y poner algun Known host
			Timeout:         _TCPTIMEOUT,
		})
		if err != nil {
			return
		}
		sshC.conn = conn
		go func() {
			for {
				if !sshC.testConnection() {
					sshC.Close()
					sshC.ErrChan <- fmt.Errorf("error de conexion")
					return
				}
				time.Sleep(_KEEPALIVETIMEOUT)
			}
		}()
	}

	return

}

func (sshC *SshClient) getSession() (sess *ssh.Session, err error) {

	var conn *ssh.Client

	conn, err = sshC.getConn()
	if err != nil {
		return
	}

	sess, err = conn.NewSession()

	return

}

func (sshC *SshClient) Close() error {
	if sshC.conn != nil {
		err := sshC.conn.Close()
		sshC.conn = nil
		return err
	}
	return nil
}

func (sshC *SshClient) ExecCommand(cmd string) (output string, err error) {

	var session *ssh.Session

	_out := new(strings.Builder)
	_err := new(strings.Builder)

	session, err = sshC.getSession()
	if err != nil {
		return
	}

	defer session.Close()

	session.Stderr = _err
	session.Stdout = _out

	err = session.Run(cmd)
	if err != nil {
		err = fmt.Errorf("command failed: %w - stderr: %s", err, _err.String())
		return
	}
	output = _out.String()
	return
}

func (sshC *SshClient) testConnection() bool {
	sess, err := sshC.conn.NewSession()
	if err != nil {
		return false
	}
	sess.Close()
	return true
}

func NewSshClient(addr, user string, privKey []byte) *SshClient {

	sshC := SshClient{addr: addr, user: user, privKey: privKey, conn: nil, ErrChan: make(chan error, 1)}

	return &sshC

}

func handleTunnel(remoteConn net.Conn, localAddr string) {

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

func (sshC *SshClient) SetupReverseTunnel(remoteAddr, localAddr string) error {
	conn, err := sshC.getConn()
	if err != nil {
		return err
	}

	listener, err := conn.Listen(_SSHTYPECONN, remoteAddr)
	if err != nil {
		return fmt.Errorf("failed to set up remote listener: %w", err)
	}

	go func() {
		for {
			remoteConn, err := listener.Accept()
			if err != nil {
				log.Printf("Tunnel accept error: %v", err)
				return
			}
			go handleTunnel(remoteConn, localAddr)
		}
	}()

	return nil

}

func (sshC *SshClient) SetupLocalTunnel(remoteAddr, localAddr string) error {

	// localServer, err := net.Listen(_SSHTYPECONN,localAddr)

	// if err != nil {
	// 	return  err
	// }
	return nil

}
