package rsync

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"slices"
	"strings"

	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/picosh/go-rsync-receiver/rsyncopts"
	"github.com/picosh/go-rsync-receiver/rsyncreceiver"
	"github.com/picosh/go-rsync-receiver/rsyncsender"
	rsyncutils "github.com/picosh/go-rsync-receiver/utils"
	"github.com/picosh/send/utils"
)

type handler struct {
	session      ssh.Session
	writeHandler utils.CopyFromClientHandler
	root         string
	recursive    bool
	ignoreTimes  bool
}

func (h *handler) List(rPath string) ([]fs.FileInfo, error) {
	isDir := false
	if rPath == "." {
		rPath = "/"
		isDir = true
	}

	list, err := h.writeHandler.List(h.session, rPath, isDir, h.recursive)
	if err != nil {
		return nil, err
	}

	var dirs []string

	var newList []fs.FileInfo

	for _, f := range list {
		if !f.IsDir() && f.Size() == 0 {
			continue
		}

		fname := f.Name()
		if strings.HasPrefix(f.Name(), "/") {
			fname = path.Join(rPath, f.Name())
		}

		if fname == "" && !f.IsDir() {
			fname = path.Base(rPath)
		}

		newFile := &utils.VirtualFile{
			FName:    fname,
			FIsDir:   f.IsDir(),
			FSize:    f.Size(),
			FModTime: f.ModTime(),
			FSys:     f.Sys(),
		}

		newList = append(newList, newFile)

		parts := strings.Split(newFile.Name(), string(os.PathSeparator))
		lastDir := newFile.Name()
		for i := 0; i < len(parts); i++ {
			lastDir, _ = path.Split(lastDir)
			if lastDir == "" {
				continue
			}

			lastDir = lastDir[:len(lastDir)-1]
			dirs = append(dirs, lastDir)
		}
	}

	for _, dir := range dirs {
		newList = append(newList, &utils.VirtualFile{
			FName:  dir,
			FIsDir: true,
		})
	}

	slices.Reverse(newList)

	onlyEmpty := true
	for _, f := range newList {
		if f.Name() != "" {
			onlyEmpty = false
		}
	}

	if len(newList) == 0 || onlyEmpty {
		return nil, errors.New("no files to send, the directory may not exist or could be empty")
	}

	return newList, nil
}

func (h *handler) Read(file *rsyncutils.SenderFile) (os.FileInfo, rsyncutils.ReaderAtCloser, error) {
	filePath := file.WPath

	if strings.HasSuffix(h.root, file.WPath) {
		filePath = h.root
	} else if !strings.HasPrefix(filePath, h.root) {
		filePath = path.Join(h.root, file.Path, file.WPath)
	}

	return h.writeHandler.Read(h.session, &utils.FileEntry{Filepath: filePath})
}

func (h *handler) Put(file *rsyncutils.ReceiverFile) (int64, error) {
	fileEntry := &utils.FileEntry{
		Filepath: path.Join("/", h.root, file.Name),
		Mode:     fs.FileMode(0600),
		Size:     file.Length,
		Mtime:    file.ModTime.Unix(),
		Atime:    file.ModTime.Unix(),
	}
	fileEntry.Reader = file.Reader

	msg, err := h.writeHandler.Write(h.session, fileEntry)
	if err != nil {
		errMsg := fmt.Sprintf("%s\r\n", err.Error())
		_, err = h.session.Stderr().Write([]byte(errMsg))
	}
	if msg != "" {
		nMsg := fmt.Sprintf("%s\r\n", msg)
		_, err = h.session.Stderr().Write([]byte(nMsg))
	}
	return 0, err
}

func (h *handler) Remove(willReceive []*rsyncutils.ReceiverFile) error {
	entries, err := h.writeHandler.List(h.session, path.Join("/", h.root), true, true)
	if err != nil {
		return err
	}

	var toDelete []string

	for _, entry := range entries {
		exists := slices.ContainsFunc(willReceive, func(rf *rsyncutils.ReceiverFile) bool {
			return rf.Name == entry.Name()
		})

		if !exists && entry.Name() != "._pico_keep_dir" {
			toDelete = append(toDelete, entry.Name())
		}
	}

	var errs []error

	for _, file := range toDelete {
		errs = append(errs, h.writeHandler.Delete(h.session, &utils.FileEntry{Filepath: path.Join("/", h.root, file)}))
		_, err = h.session.Stderr().Write([]byte(fmt.Sprintf("deleting %s\r\n", file)))
		errs = append(errs, err)
	}

	return errors.Join(errs...)
}

func Middleware(writeHandler utils.CopyFromClientHandler) wish.Middleware {
	return func(sshHandler ssh.Handler) ssh.Handler {
		return func(session ssh.Session) {
			cmd := session.Command()
			if len(cmd) == 0 || cmd[0] != "rsync" {
				sshHandler(session)
				return
			}

			logger := writeHandler.GetLogger(session).With(
				"rsync", true,
				"cmd", cmd,
			)

			defer func() {
				if r := recover(); r != nil {
					logger.Error("error running rsync middleware", "err", r)
					_, _ = session.Stderr().Write([]byte("error running rsync middleware, check the flags you are using\r\n"))
				}
			}()

			cmdFlags := session.Command()

			optsCtx, err := rsyncopts.ParseArguments(cmdFlags[1:], true)
			if err != nil {
				_, _ = session.Stderr().Write([]byte(fmt.Sprintf("error parsing rsync arguments: %s\r\n", err.Error())))
				return
			}

			if optsCtx.Options.Compress() {
				_, _ = session.Stderr().Write([]byte("compression is currently unsupported\r\n"))
				return
			}

			if len(optsCtx.RemainingArgs) != 2 {
				_, _ = session.Stderr().Write([]byte("missing source and destination arguments\r\n"))
				return
			}

			root := strings.TrimPrefix(optsCtx.RemainingArgs[len(optsCtx.RemainingArgs)-1], "/")
			if root == "" {
				root = "/"
			}

			fileHandler := &handler{
				session:      session,
				writeHandler: writeHandler,
				root:         root,
				recursive:    optsCtx.Options.Recurse(),
				ignoreTimes:  !optsCtx.Options.PreserveMTimes(),
			}

			for _, arg := range cmd {
				if arg == "--sender" {
					if err := rsyncsender.ClientRun(logger, optsCtx.Options, session, fileHandler, []string{fileHandler.root}, true); err != nil {
						logger.Error("error running rsync sender", "err", err)
					}
					return
				}
			}

			err = rsyncreceiver.ClientRun(logger, optsCtx.Options, session, fileHandler, []string{fileHandler.root}, true)
			if err != nil {
				logger.Error("error running rsync receiver", "err", err)
				return
			}
		}
	}
}
