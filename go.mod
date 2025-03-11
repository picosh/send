module github.com/picosh/send

go 1.24

toolchain go1.24.0

// replace github.com/picosh/go-rsync-receiver => ../go-rsync-receiver

replace github.com/picosh/pico => ../pico

require (
	github.com/charmbracelet/ssh v0.0.0-20250213143314-8712ec3ff3ef
	github.com/charmbracelet/wish v1.4.6
	github.com/matryer/is v1.4.1
	github.com/picosh/go-rsync-receiver v0.0.0-20250304201040-fcc11dd22d79
	github.com/picosh/pico v0.0.0-00010101000000-000000000000
	github.com/pkg/sftp v1.13.7
)

require (
	github.com/anmitsu/go-shlex v0.0.0-20200514113438-38f4b401e2be // indirect
	github.com/antoniomika/syncmap v1.0.0 // indirect
	github.com/aymanbagabas/go-osc52/v2 v2.0.1 // indirect
	github.com/charmbracelet/bubbletea v1.3.4 // indirect
	github.com/charmbracelet/keygen v0.5.1 // indirect
	github.com/charmbracelet/lipgloss v1.0.0 // indirect
	github.com/charmbracelet/log v0.4.0 // indirect
	github.com/charmbracelet/x/ansi v0.8.0 // indirect
	github.com/charmbracelet/x/conpty v0.1.0 // indirect
	github.com/charmbracelet/x/errors v0.0.0-20250226164017-59292a315e58 // indirect
	github.com/charmbracelet/x/term v0.2.1 // indirect
	github.com/charmbracelet/x/termios v0.1.1 // indirect
	github.com/creack/pty v1.1.24 // indirect
	github.com/erikgeiser/coninput v0.0.0-20211004153227-1c3628e74d0f // indirect
	github.com/go-logfmt/logfmt v0.6.0 // indirect
	github.com/kr/fs v0.1.0 // indirect
	github.com/lucasb-eyer/go-colorful v1.2.0 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-localereader v0.0.1 // indirect
	github.com/mattn/go-runewidth v0.0.16 // indirect
	github.com/mmcloughlin/md4 v0.1.2 // indirect
	github.com/muesli/ansi v0.0.0-20230316100256-276c6243b2f6 // indirect
	github.com/muesli/cancelreader v0.2.2 // indirect
	github.com/muesli/termenv v0.16.0 // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	golang.org/x/crypto v0.36.0 // indirect
	golang.org/x/exp v0.0.0-20250218142911-aa4b98e5adaa // indirect
	golang.org/x/sync v0.12.0 // indirect
	golang.org/x/sys v0.31.0 // indirect
	golang.org/x/text v0.23.0 // indirect
)
