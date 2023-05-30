package sender

import (
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"net"
	"time"
)

func NewTCPSender(host, port string, duration time.Duration) (Interface, error) {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%s", host, port), duration)
	if err != nil {
		return nil, errors.Wrap(err, "unable to establish a TCP connection")
	}

	return &tcpSender{
			conn: conn,
		}, errors.Wrap(
			conn.SetDeadline(time.Now().Add(duration*2)),
			"unable to set connection deadline",
		)
}

var _ Interface = (*tcpSender)(nil)

type tcpSender struct {
	conn net.Conn
}

func (t *tcpSender) Close() error {
	return t.conn.Close()
}

func (t *tcpSender) Send(msg []byte) error {
	_, err := io.Copy(t.conn, bytes.NewReader(msg))
	return err
}
