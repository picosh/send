package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/picosh/send/protocols"
	"github.com/picosh/send/utils"
)

type handler struct {
}

func (h *handler) GetLogger() *slog.Logger {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	return logger
}

func (h *handler) Delete(session ssh.Session, file *utils.FileEntry) error {
	str := fmt.Sprintf("Deleted file: %+v from session: %+v", file, session)
	log.Print(str)
	return nil
}

func (h *handler) Write(session ssh.Session, file *utils.FileEntry) (string, error) {
	str := fmt.Sprintf("Received file: %+v from session: %+v", file, session)
	log.Print(str)
	return str, nil
}

func (h *handler) Validate(session ssh.Session) error {
	log.Printf("Received validate from session: %+v", session)

	return nil
}

func (h *handler) Read(session ssh.Session, entry *utils.FileEntry) (os.FileInfo, utils.ReaderAtCloser, error) {
	log.Printf("Received validate from session: %+v", session)

	data := strings.NewReader("lorem ipsum dolor")

	return &utils.VirtualFile{
		FName:    "test",
		FIsDir:   false,
		FSize:    data.Size(),
		FModTime: time.Now(),
	}, utils.NopReaderAtCloser(data), nil
}

func (h *handler) List(session ssh.Session, fpath string, isDir bool, recursive bool) ([]os.FileInfo, error) {
	return nil, nil
}

func main() {
	h := &handler{}

	s, err := wish.NewServer(
		wish.WithAddress("localhost:9000"),
		wish.WithHostKeyPath("ssh_data/term_info_ed25519"),
		protocols.Middleware(h),
	)

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Starting ssh server on 9000")

	log.Fatal(s.ListenAndServe())
}
