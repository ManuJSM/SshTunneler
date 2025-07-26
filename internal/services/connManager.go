package services

import (
	"fmt"
	"time"

	"golang.org/x/crypto/ssh"
)

type SSHConnection interface {
	Connect() error
	Close() error
	IsAlive() bool
	Client() (*ssh.Client, error)
}

type SSHConnectionImpl struct {
	addr    string
	user    string
	privKey []byte
	conn    *ssh.Client
}

const _SSHTYPECONN = "tcp"
const _TCPTIMEOUT = 1 * time.Second

func NewSSHConnection(addr, user string, privKey []byte) SSHConnection {

	return &SSHConnectionImpl{addr: addr, user: user, privKey: privKey, conn: nil}

}
func getAuthMethod(privKey []byte) (ssh.AuthMethod, error) {

	signer, err := ssh.ParsePrivateKey(privKey)
	if err != nil {
		return nil, err
	}
	return ssh.PublicKeys(signer), nil
}

func (sshC *SSHConnectionImpl) Connect() (err error) {

	// se comprueba si ya esta establecida la conexion
	if sshC.conn == nil {

		var authMethod ssh.AuthMethod

		authMethod, err = getAuthMethod(sshC.privKey)
		if err != nil {
			return
		}

		sshC.conn, err = ssh.Dial(_SSHTYPECONN, sshC.addr, &ssh.ClientConfig{
			User:            sshC.user,
			Auth:            []ssh.AuthMethod{authMethod},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(), // ⚠️ Conveniente cambiarlo y poner algun Known host
			Timeout:         _TCPTIMEOUT,
		})
	}

	return

}

func (sshC *SSHConnectionImpl) Close() (err error) {
	if sshC.conn != nil {
		err = sshC.conn.Close()
		sshC.conn = nil
	}
	return
}

func (sshC *SSHConnectionImpl) IsAlive() bool {
	sess, err := sshC.conn.NewSession()
	if err != nil {
		return false
	}
	sess.Close()
	return true
}

func (sshC *SSHConnectionImpl) Client() (*ssh.Client, error) {
	if !sshC.IsAlive() {
		return nil, fmt.Errorf("cliente esta Riperino capuchino ")
	}
	return sshC.conn, nil
}
