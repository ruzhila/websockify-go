package websockifygo

import (
	"crypto/tls"
	"io"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
)

type WSproxy struct {
	URL     string
	Target  string
	KeyPem  string
	CertPem string
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (s *WSproxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != s.URL || r.Method != http.MethodGet {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	defer conn.Close()
	target, err := net.Dial("tcp", s.Target)
	if err != nil {
		http.Error(w, "Bad Gateway", http.StatusBadGateway)
		return
	}
	defer target.Close()
	errChan := make(chan error, 2)
	go func() {
		for {
			messageType, message, err := conn.ReadMessage()
			if err != nil {
				errChan <- err
				return
			}
			if messageType != websocket.BinaryMessage && messageType != websocket.TextMessage {
				continue
			}
			if _, err = target.Write(message); err != nil {
				errChan <- err
				return
			}
		}
	}()
	go func() {
		buf := make([]byte, 1024)
		var n int
		for {
			if n, err = target.Read(buf); err != nil {
				errChan <- err
				return
			}
			if err = conn.WriteMessage(websocket.BinaryMessage, buf[:n]); err != nil {
				errChan <- err
				return
			}
		}
	}()
	log.Println("proxying", s.URL, "to", s.Target, "remote", conn.RemoteAddr())
	err = <-errChan
	if err != io.EOF {
		log.Println("proxy error:", err)
	}
}

func (prx *WSproxy) Serve(addr string) error {
	proxyNetListner, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	httpsrv := &http.Server{
		Handler: prx,
	}
	if strings.HasSuffix(addr, "443") && prx.CertPem != "" && prx.KeyPem != "" {
		httpsrv.TLSConfig = &tls.Config{
			MinVersion:       tls.VersionTLS12,
			CurvePreferences: []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		}
		httpsrv.TLSNextProto = make(map[string]func(*http.Server, *tls.Conn, http.Handler))
		return httpsrv.ServeTLS(proxyNetListner, prx.CertPem, prx.KeyPem)
	} else {
		return httpsrv.Serve(proxyNetListner)
	}
}
