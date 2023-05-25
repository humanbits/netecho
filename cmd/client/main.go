package main

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
	"io"
	"net"
	"time"
)

type config struct {
	TargetHostname   string        `default:"localhost"`
	TargetPort       int           `default:"8080"`
	MessageSizeBytes int64         `default:"10000000"`
	SleepDuration    time.Duration `default:"1s"`
	Timeout          time.Duration `default:"100ms"`
	LogLevel         string        `default:"info"`
}

func main() {
	var c config
	if err := envconfig.Process("", &c); err != nil {
		log.Fatalf("config error: %s", err.Error())
	}

	level, _ := log.ParseLevel(c.LogLevel)
	log.SetLevel(level)

	for {
		conn, err := net.DialTimeout(
			"tcp",
			fmt.Sprintf("%s:%d", c.TargetHostname, c.TargetPort),
			c.Timeout,
		)

		if err != nil {
			log.Errorf("unable to establish connection: %s", err.Error())
		} else {
			if err := conn.SetDeadline(time.Now().Add(c.Timeout)); err != nil {
				log.Fatalf("error setting connection timeout: %v:", err)
			}
			originalMessage, err := io.ReadAll(io.LimitReader(rand.Reader, c.MessageSizeBytes))
			log.Printf(
				"successfully established connection, sending %d random bytes (%s)...",
				c.MessageSizeBytes,
				base64.StdEncoding.EncodeToString(originalMessage[:10]),
			)
			if err != nil {
				log.Error(err.Error())
			}

			_, err = io.Copy(conn, bytes.NewReader(originalMessage))
			if err != nil {
				log.Errorf("error sending bytes: %s", err.Error())
			} else {
				log.Printf("OK, closing...")
			}
			if err := conn.Close(); err != nil {
				log.Errorf("error closing the connection: %s", err.Error())
			} else {
				log.Info("OK")
			}
		}
		log.Printf("sleeping %v...", c.SleepDuration)
		time.Sleep(c.SleepDuration)
	}

}
