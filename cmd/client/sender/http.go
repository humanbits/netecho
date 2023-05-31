package sender

import (
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"net/http"
	"time"
)

func NewHTTPSender(host, port string, duration time.Duration) (Interface, error) {
	return &httpSender{
		host:     host,
		port:     port,
		duration: duration,
	}, nil
}

type httpSender struct {
	host, port string
	duration   time.Duration
}

func (h *httpSender) Send(msg []byte) error {
	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("http://%s:%s", h.host, h.port),
		bytes.NewReader(msg),
	)
	req.Header.Set("Keep-Alive", fmt.Sprintf("max=%d", int(h.duration.Seconds())))
	if err != nil {
		return errors.Wrap(err, "error building request")
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "error sending request")
	}
	if err := res.Body.Close(); err != nil {
		return errors.Wrap(err, "unable to close the body")
	}
	if res.StatusCode != http.StatusOK {
		return errors.Errorf("response code is %d not 200", res.StatusCode)
	}
	return nil
}

func (h *httpSender) Close() error {
	return nil
}
