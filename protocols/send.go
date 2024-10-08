package protocols

import (
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/picosh/send/auth"
	"github.com/picosh/send/pipe"
	"github.com/picosh/send/protocols/rsync"
	"github.com/picosh/send/protocols/scp"
	"github.com/picosh/send/protocols/sftp"
	"github.com/picosh/send/utils"
)

func Middleware(writeHandler utils.CopyFromClientHandler) ssh.Option {
	return func(server *ssh.Server) error {
		err := wish.WithMiddleware(pipe.Middleware(writeHandler, ""), scp.Middleware(writeHandler), rsync.Middleware(writeHandler), auth.Middleware(writeHandler))(server)
		if err != nil {
			return err
		}

		return sftp.SSHOption(writeHandler)(server)
	}
}
