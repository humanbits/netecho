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
	ConnDuration     time.Duration `default:"100ms"`
	SleepDuration    time.Duration `default:"10ms"`
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
			c.ConnDuration,
		)

		if err != nil {
			log.Errorf("unable to establish connection: %s", err.Error())
		} else {
			if err := conn.SetDeadline(time.Now().Add(c.ConnDuration * 2)); err != nil {
				log.Fatalf("error setting connection timeout: %v:", err)
			}
			originalMessage, err := io.ReadAll(io.LimitReader(rand.Reader, c.MessageSizeBytes))
			log.Printf(
				"successfully established new connection, it will be sending %d random bytes (%s...) up to %d times over it",
				c.MessageSizeBytes,
				base64.StdEncoding.EncodeToString(originalMessage[:10]),
				c.ConnDuration/c.SleepDuration,
			)
			if err != nil {
				log.Error(err.Error())
			}

			start := time.Now()
			i := 0
			var errors []error
			for ; start.Add(c.ConnDuration).After(time.Now()); i++ {
				log.Infof("sending message %d...", i)
				_, err = io.Copy(conn, bytes.NewReader(originalMessage))
				if err != nil {
					errors = append(errors, err)
					log.Errorf("error sending bytes: %s", err.Error())
				}
				time.Sleep(c.SleepDuration)
			}

			log.Printf("successfully sent %d of %d messages, closing the connection...", len(errors), i)

			if err := conn.Close(); err != nil {
				log.Errorf("error closing the connection: %s", err.Error())
			} else {
				log.Info("successfully closed connection")
			}
		}
	}
}
