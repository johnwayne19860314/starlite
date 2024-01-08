package features

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"

	"github.startlite.cn/itapp/startlite/pkg/lines/appx"
	"github.startlite.cn/itapp/startlite/pkg/lines/errorx"
	"github.startlite.cn/itapp/startlite/pkg/lines/featurex"
	"github.startlite.cn/itapp/startlite/pkg/lines/utilx/textx"
	"github.startlite.cn/itapp/startlite/pkg/servicex/types"
)

type Sftp struct {
	types.SftpConfig

	*sftp.Client
	sshConfig *ssh.ClientConfig
}

func MustNewSftp(appCtx appx.AppContext, cl *featurex.ConfigLoader) *Sftp {
	sftpClient := Sftp{}
	cl.Load(&sftpClient)

	sftpClient.PrivateKey = strings.ReplaceAll(sftpClient.PrivateKey, "\\n", "\n")

	config := &ssh.ClientConfig{
		User:            sftpClient.User,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	if !textx.Blank(sftpClient.Password) {
		config.Auth = append(config.Auth, ssh.Password(sftpClient.Password))
	}
	if !textx.Blank(sftpClient.PrivateKey) {
		signer, err := ssh.ParsePrivateKey([]byte(sftpClient.PrivateKey))
		if err != nil {
			appCtx.Fatal("unable to parse private key: %v", err)
		}
		config.Auth = append(config.Auth, ssh.PublicKeys(signer))
	}

	sftpClient.sshConfig = config

	err := sftpClient.SetNewClient()
	if err != nil {
		appCtx.Fatal(err.Error())
	}

	return &sftpClient
}

func (client *Sftp) SetNewClient() error {
	if client.Client != nil {
		_ = client.Client.Close()
	}

	sshClient, err := ssh.Dial("tcp", fmt.Sprintf("%s:%s", client.Host, client.Port), client.sshConfig)
	if err != nil {
		return errorx.Errorf("unable to connect: %v", err)
	}
	// defer sshClient.Close()

	client.Client, err = sftp.NewClient(sshClient)
	if err != nil {
		return errorx.Errorf("unable to new client: %v", err)
	}

	return nil
}

func (client *Sftp) TryDo(f func() error) error {
	err := f()
	if err == nil {
		return nil
	}

	if !strings.Contains(err.Error(), "connection lost") {
		return err
	}

	err = client.SetNewClient()
	if err != nil {
		return err
	}

	return f()
}

func (client *Sftp) Upload(remotePath string, reader io.Reader) error {
	return client.TryDo(func() error {
		dstFile, err := client.OpenFile(remotePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC)
		if err != nil {
			return errorx.WithStack(err)
		}
		defer dstFile.Close()

		_, err = io.Copy(dstFile, reader)
		if err != nil {
			return errorx.WithStack(err)
		}

		return nil
	})
}
