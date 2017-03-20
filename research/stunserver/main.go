package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"github.com/jessevdk/go-flags"
	"github.com/osrtss/icego/stun"
)

var opts struct {
	Host string `long:"host" default:"0.0.0.0" description:"IP to bind to"`
	Port uint16 `long:"port" default:"2202" description:"UDP port to bind to"`
	File string `long:"file" default:"" description:"dump received data to a dump file"`
}

func main() {
	_, err := flags.Parse(&opts)
	if err != nil {
		if !strings.Contains(err.Error(), "Usage") {
			log.Printf("error: %v\n", err.Error())
			os.Exit(1)
		} else {
			log.Printf("%v\n", err.Error())
			os.Exit(0)
		}
	}

	srv := stun.NewServer(nil)
	l, err := net.ListenPacket("udp", fmt.Sprintf("%s:%d", opts.Host, opts.Port))
	if err != nil {
		log.Fatal("Listen error: ", err)
	}
	defer l.Close()

	log.Printf(">> Starting stundump, listening at %v:%v...", opts.Host, opts.Port)
	srv.ServePacket(l)
}
