package auth

import (
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/picosh/send/utils"
)

func Middleware(writeHandler utils.CopyFromClientHandler) wish.Middleware {
	return func(sshHandler ssh.Handler) ssh.Handler {
		return func(session ssh.Session) {
			defer func() {
				if r := recover(); r != nil {
					writeHandler.GetLogger().Error("error running auth middleware", "err", r)
				}
			}()

			err := writeHandler.Validate(session)
			if err != nil {
				utils.ErrorHandler(session, err)
				return
			}

			sshHandler(session)
		}
	}
}
