package svrkit

import (
	"fmt"
	"net"
	"os"
	u "os/user"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
)

// SSHClient SSH 客户端，这里的方法主要用于跨机器命令操作
type SSHClient struct {
	*ssh.Client
}

// RunCommand ssh 远程命令执行 返回执行结果输出
func (client *SSHClient) RunCommand(cmd string) ([]byte, error) {
	sess, err := client.NewSession()
	if err != nil {
		return nil, err
	}
	defer sess.Close()

	return sess.CombinedOutput(cmd)
}

// WriteFile ssh 远程文件写入 读取可以通过执行 cat 命令实现
func (client *SSHClient) WriteFile(targetPath string, data []byte) error {
	filename := filepath.Base(targetPath)
	dirname := strings.Replace(filepath.Dir(targetPath), "\\", "/", -1)

	sess, err := client.NewSession()
	if err != nil {
		return err
	}
	defer sess.Close()

	go func() {
		w, _ := sess.StdinPipe()
		fmt.Fprintln(w, "C0644", len(data), filename)
		w.Write(data)
		fmt.Fprint(w, "\x00")
		w.Close()
	}()

	return sess.Run(fmt.Sprintf("/usr/bin/scp -qrt %s", dirname))
}

// NewSSHClientConnectByPass 使用用户名密码方式创建 ssh 连接
func NewSSHClientConnectByPass(user, pass, addr string) (*SSHClient, error) {
	client, err := ssh.Dial("tcp", addr, &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(pass),
		},
		Timeout: 30 * time.Second,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	})
	if err != nil {
		return nil, err
	}

	return &SSHClient{client}, nil
}

// NewSSHClientConnectByKey 使用密钥方式创建 ssh 连接
func NewSSHClientConnectByKey(user, privateKeyFile, addr string) (*SSHClient, error) {
	if privateKeyFile[:2] == "~/" {
		usr, err := u.Current()
		if err != nil {
			return nil, err
		}
		privateKeyFile = filepath.Join(usr.HomeDir, privateKeyFile[2:])
	}

	key, err := os.ReadFile(privateKeyFile)
	if err != nil {
		return nil, err
	}

	// Create the Signer for this private key.
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, err
	}

	client, err := ssh.Dial("tcp", addr, &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		Timeout: 30 * time.Second,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	})
	if err != nil {
		return nil, err
	}

	return &SSHClient{client}, nil
}
