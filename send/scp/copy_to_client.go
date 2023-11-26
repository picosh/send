package scp

import (
	"errors"

	"github.com/charmbracelet/ssh"
	"github.com/picosh/send/send/utils"
)

func copyToClient(session ssh.Session, info Info, handler utils.CopyFromClientHandler) error {
	return errors.New("unsupported, use rsync or sftp")
}
