package sftp

import (
	"errors"
	"io"

	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/picosh/send/utils"
	"github.com/pkg/sftp"
)

func SSHOption(writeHandler utils.CopyFromClientHandler) ssh.Option {
	return func(server *ssh.Server) error {
		if server.SubsystemHandlers == nil {
			server.SubsystemHandlers = map[string]ssh.SubsystemHandler{}
		}

		server.SubsystemHandlers["sftp"] = SubsystemHandler(writeHandler)
		return nil
	}
}

func SubsystemHandler(writeHandler utils.CopyFromClientHandler) ssh.SubsystemHandler {
	return func(session ssh.Session) {
		logger := writeHandler.GetLogger(session).With(
			"sftp", true,
		)

		defer func() {
			if r := recover(); r != nil {
				logger.Error("error running sftp middleware", "err", r)
				wish.Println(session, "error running sftp middleware, check the flags you are using")
			}
		}()

		err := writeHandler.Validate(session)
		if err != nil {
			wish.Errorln(session, err)
			return
		}

		handler := &handlererr{
			Handler: &handler{
				session:      session,
				writeHandler: writeHandler,
			},
		}

		handlers := sftp.Handlers{
			FilePut:  handler,
			FileList: handler,
			FileGet:  handler,
			FileCmd:  handler,
		}

		requestServer := sftp.NewRequestServer(session, handlers)

		err = requestServer.Serve()
		if err != nil && !errors.Is(err, io.EOF) {
			wish.Errorln(session, err)
			logger.Error("Error serving sftp subsystem", "err", err)
		}
	}
}
