package main

import (
	"flag"
	"fmt"
	"github.com/kelseyhightower/envconfig"
	"io"
	"log"
	"net"
)

type config struct {
	Port int `default:"8080" desc:"Port to listen on"`
}

func main() {
	showUsage := flag.Bool("h", false, "show usage description")
	flag.Parse()

	var c config
	if err := envconfig.Process("", &c); err != nil {
		log.Fatalf("error reading config: %s", err.Error())
	}
	if *showUsage {
		err := envconfig.Usage("", &c)
		if err != nil {
			log.Fatalf("could not show usage: %s", err.Error())
		}
		return
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", c.Port))
	if err != nil {
		log.Fatalf("could not listen on the port: %s", err.Error())
	}
	log.Printf("started listening on %d", c.Port)

	for {
		con, err := lis.Accept()
		if err != nil {
			log.Printf("could not accept connection: %s", err.Error())
			continue
		}
		log.Printf("accepted new connection")

		n, err := io.Copy(con, con)
		if err != nil {
			log.Printf("error echoing input: %s", err.Error())
			continue
		}
		log.Printf("successfully echoed %d bytes back", n)
	}
}
