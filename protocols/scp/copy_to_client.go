package scp

import (
	"errors"

	"github.com/picosh/pico/pssh"
	"github.com/picosh/send/utils"
)

func copyToClient(session *pssh.SSHServerConnSession, info Info, handler utils.CopyFromClientHandler) error {
	return errors.New("unsupported, use rsync or sftp")
}
