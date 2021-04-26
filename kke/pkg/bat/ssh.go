package bat

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

type Client struct {
	SSHClient  *ssh.Client
	SSHSession *ssh.Session
	SFTPClient *sftp.Client
}

func New(host Host) (*Client, error) {
	sshClient, sshSession, err := connect(host)
	if err != nil {
		return nil, err
	}
	sftpClient, err := sftp.NewClient(sshClient)
	if err != nil {
		return nil, err
	}

	return &Client{
		SSHClient:  sshClient,
		SSHSession: sshSession,
		SFTPClient: sftpClient,
	}, nil
}

func connect(host Host) (*ssh.Client, *ssh.Session, error) {
	var authMethod ssh.AuthMethod
	if host.Cred == "" {
		home, _ := os.UserHomeDir()
		host.Cred = path.Join(home, ".ssh/id_rsa")
	}
	if pemBytes, err := ioutil.ReadFile(host.Cred); err != nil {
		authMethod = ssh.Password(host.Cred)
	} else {
		signer, err := ssh.ParsePrivateKey(pemBytes)
		if err != nil {
			return nil, nil, err
		}
		authMethod = ssh.PublicKeys(signer)
	}

	sshConfig := &ssh.ClientConfig{
		User: host.User,
		Auth: []ssh.AuthMethod{authMethod},
	}

	sshConfig.HostKeyCallback = ssh.InsecureIgnoreHostKey()

	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", host.Addr, host.Port), sshConfig)
	if err != nil {
		return nil, nil, err
	}

	session, err := client.NewSession()
	if err != nil {
		client.Close()
		return nil, nil, err
	}

	return client, session, nil
}

func downloadFile(sftpClient *sftp.Client, remotePath string, localPath string) error {
	remoteFile, err := sftpClient.Open(remotePath)
	if err != nil {
		return err

	}
	defer remoteFile.Close()

	_ = os.MkdirAll(path.Dir(localPath), os.ModePerm)
	localFile, err := os.Create(localPath)
	if err != nil {
		localFile, err = os.Create(path.Join(localPath, path.Base(remotePath)))
		if err != nil {
			return err
		}
	}
	defer localFile.Close()

	_, err = remoteFile.WriteTo(localFile)
	return err
}

func downloadDir(sftpClient *sftp.Client, remotePath string, localPath string) error {
	remoteFiles, err := sftpClient.ReadDir(remotePath)
	if err != nil {
		return err
	}

	ff, err := os.Stat(localPath)
	if err == nil && ff.IsDir() {
		if path.Base(remotePath) != path.Base(localPath) {
			localPath = path.Join(localPath, path.Base(remotePath))
		}
	}
	_ = os.MkdirAll(localPath, os.ModePerm)

	for _, backupDir := range remoteFiles {
		remoteFilePath := path.Join(remotePath, backupDir.Name())
		localFilePath := path.Join(localPath, backupDir.Name())
		if backupDir.IsDir() {
			err = os.MkdirAll(localFilePath, os.ModePerm)
			if err != nil {
				return err
			}
			err = downloadDir(sftpClient, remoteFilePath, localFilePath)
			if err != nil {
				return err
			}
		} else {
			err = downloadFile(sftpClient, remoteFilePath, localFilePath)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func uploadFile(sftpClient *sftp.Client, localPath string, remotePath string) error {
	localFile, err := os.Open(localPath)
	if err != nil {
		return err

	}
	defer localFile.Close()

	remoteFile, err := sftpClient.Create(remotePath)
	if err != nil {
		remoteFile, err = sftpClient.Create(path.Join(remotePath, path.Base(localPath)))
		if err != nil {
			return err
		}
	}
	defer remoteFile.Close()

	ff, err := ioutil.ReadAll(localFile)
	if err != nil {
		return err

	}

	_, err = remoteFile.Write(ff)
	return err
}

func uploadDir(sftpClient *sftp.Client, localPath string, remotePath string) error {
	if sftpClient.Mkdir(remotePath) != nil { // remotePath is already exsit in remote host
		if path.Base(localPath) != path.Base(remotePath) {
			remotePath = path.Join(remotePath, path.Base(localPath))
		}
		_ = sftpClient.MkdirAll(remotePath)
	}

	localFiles, err := ioutil.ReadDir(localPath)
	if err != nil {
		return err
	}

	for _, backupDir := range localFiles {
		localFilePath := path.Join(localPath, backupDir.Name())
		remoteFilePath := path.Join(remotePath, backupDir.Name())
		if backupDir.IsDir() {
			err = sftpClient.MkdirAll(remoteFilePath)
			if err != nil {
				return err
			}
			err = uploadDir(sftpClient, localFilePath, remoteFilePath)
			if err != nil {
				return err
			}
		} else {
			err = uploadFile(sftpClient, path.Join(localPath, backupDir.Name()), remotePath)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (c *Client) Execute(command string) ([]byte, error) {
	return c.SSHSession.CombinedOutput(command)
}

func (c *Client) Script(scriptPath string) ([]byte, error) {
	return c.SSHSession.CombinedOutput(fmt.Sprintf("sh %s", scriptPath))
}

func (c *Client) Template(remotePath, oldStr, newStr string) ([]byte, error) {
	return c.SSHSession.CombinedOutput(fmt.Sprintf("sed -i ^s/%s^%s^g %s", oldStr, newStr, remotePath))
}

func (c *Client) Shell(command string) error {
	// Set IO
	c.SSHSession.Stdout = os.Stdout
	c.SSHSession.Stderr = os.Stderr
	in, _ := c.SSHSession.StdinPipe()

	// Set up terminal modes
	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	// Request pseudo terminal
	if err := c.SSHSession.RequestPty("xterm", 80, 40, modes); err != nil {
		return err
	}

	// Start remote shell
	if err := c.SSHSession.Shell(); err != nil {
		return err
	}

	if command != "" {
		fmt.Fprintln(in, command)
	}
	// Accepting commands
	for {
		reader := bufio.NewReader(os.Stdin)
		str, _ := reader.ReadString('\n')
		fmt.Fprint(in, str)
	}
}

func (c *Client) Pull(remotePath, localPath string) error {
	if localPath == "" {
		localPath = remotePath
	}
	localPath, err := filepath.Abs(localPath)
	if err != nil {
		return err
	}

	_, err = c.SFTPClient.ReadDir(remotePath)
	if err != nil {
		return downloadFile(c.SFTPClient, remotePath, localPath)
	}
	return downloadDir(c.SFTPClient, remotePath, localPath)
}

func (c *Client) Push(localPath, remotePath string) error {
	ff, err := os.Stat(localPath)
	if err != nil {
		return err
	}
	localPath, err = filepath.Abs(localPath)
	if err != nil {
		return err
	}
	if remotePath == "" {
		remotePath = localPath
	}

	if ff.IsDir() {
		return uploadDir(c.SFTPClient, localPath, remotePath)
	}

	return uploadFile(c.SFTPClient, localPath, remotePath)
}

func (c *Client) Close() {
	c.SFTPClient.Close()
	c.SSHSession.Close()
	c.SSHClient.Close()
}
