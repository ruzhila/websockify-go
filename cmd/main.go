package main

import (
	"flag"
	"log"
	"strings"

	websockifygo "github.com/ruzhila/websockify-go"
)

func main() {
	var addr string
	var target string
	var keyPem string
	var certPem string
	var url string

	flag.StringVar(&addr, "addr", ":8080", "address to listen on")
	flag.StringVar(&target, "target", "localhost:5900", "target address")
	flag.StringVar(&keyPem, "key", "", "SSL key.pem")
	flag.StringVar(&certPem, "cert", "", "SSL cert.pem")
	flag.StringVar(&url, "url", "", "url path to proxy, e.g. /vnc")

	flag.Parse()

	if !strings.HasPrefix(url, "/") {
		url = "/" + url
	}

	wsp := &websockifygo.WSproxy{
		URL:     url,
		Target:  target,
		KeyPem:  keyPem,
		CertPem: certPem,
	}
	log.Println("Proxying", url, "to", target, "on", addr)
	wsp.Serve(addr)
}
