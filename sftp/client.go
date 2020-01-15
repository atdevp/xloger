package sftp

import (
	"fmt"
	"net"

	"github.com/atdevp/devlib/file"
	sftp_client "github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

type SSHConfig struct {
	Host        string
	Port        int
	User        string
	RSAFilePath string
	Passwd      string
	SSHSession  *ssh.Client
}

func (this *SSHConfig) getSSHClientConfig() (*ssh.ClientConfig, error) {

	config := &ssh.ClientConfig{
		User: this.User,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	if this.Passwd != "" {
		config.Auth = []ssh.AuthMethod{
			ssh.Password(this.Passwd),
		}
		return config, nil
	}

	var p = this.RSAFilePath
	pkey, err := file.ToBytes(p)
	if err != nil {
		return nil, err
	}

	sign, err := ssh.ParsePrivateKey(pkey)
	if err != nil {
		return nil, err
	}

	config.Auth = []ssh.AuthMethod{
		ssh.PublicKeys(sign),
	}

	return config, nil
}

func (this *SSHConfig) GetSocket() string {
	return fmt.Sprintf("%s:%d", this.Host, this.Port)
}

func (this *SSHConfig) NewSFTPClient() (*sftp_client.Client, error) {
	sock := this.GetSocket()

	config, err := this.getSSHClientConfig()
	if err != nil {
		return nil, err
	}

	session, err := ssh.Dial("tcp", sock, config)
	if err != nil {
		return nil, err
	}
	this.SSHSession = session

	client, err := sftp_client.NewClient(session)
	if err != nil {
		return nil, err
	}
	return client, nil

}
