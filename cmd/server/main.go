package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
	"io"
	"net"
	"net/http"
)

type config struct {
	Port      int    `default:"8080" desc:"Port to listen on"`
	Transport string `default:"tcp"`
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

	port := c.Port
	serverFunc, ok := serverFuncs[c.Transport]
	if !ok {
		log.Fatalf("unknown transport %q, allowed values are tcp, http", c.Transport)
	}
	serverFunc(port)
}

var serverFuncs = map[string]func(port int){
	"tcp":  runTCPServer,
	"http": runHTTPServer,
}

func runTCPServer(port int) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("could not listen on the port: %s", err.Error())
	}
	log.Printf("started listening on %d", port)

	for {
		con, err := lis.Accept()
		if err != nil {
			log.Errorf("could not accept connection: %s", err.Error())
			continue
		}
		log.Printf("accepted new connection")

		n, err := io.Copy(base64.NewEncoder(base64.StdEncoding, io.Discard), con)
		if err != nil {
			log.Errorf("error echoing input: %s", err.Error())
			continue
		}
		log.Printf("successfully read %d bytes", n)
	}
}

func runHTTPServer(port int) {
	log.Printf("starting HTTP server on port %d...", port)
	handler := http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		n, err := io.Copy(io.Discard, request.Body)
		if err != nil {
			log.Errorf("unable to read request body: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			log.Infof("read %d bytes", n)
		}
	})
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), handler)
	if err != nil {
		log.Fatalf("error running http server: %v", err)
	}
}
