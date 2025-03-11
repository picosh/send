module github.com/picosh/send

go 1.24

toolchain go1.24.0

// replace github.com/picosh/go-rsync-receiver => ../go-rsync-receiver

replace github.com/picosh/pico => ../pico

require (
	github.com/matryer/is v1.4.1
	github.com/picosh/go-rsync-receiver v0.0.0-20250304201040-fcc11dd22d79
	github.com/picosh/pico v0.0.0-00010101000000-000000000000
	github.com/pkg/sftp v1.13.7
)

require (
	github.com/antoniomika/syncmap v1.0.0 // indirect
	github.com/kr/fs v0.1.0 // indirect
	github.com/mmcloughlin/md4 v0.1.2 // indirect
	github.com/stretchr/testify v1.9.0 // indirect
	golang.org/x/crypto v0.36.0 // indirect
	golang.org/x/sync v0.12.0 // indirect
	golang.org/x/sys v0.31.0 // indirect
)
