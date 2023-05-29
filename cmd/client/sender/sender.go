package sender

import "io"

type Interface interface {
	Send(msg []byte) error
	io.Closer
}
