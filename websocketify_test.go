package websockifygo

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

func TestProxy(t *testing.T) {
	// create a example echo server
	echoServer, err := net.Listen("tcp", "127.0.0.1:12340")
	if err != nil {
		t.Fatal(err)
	}

	go func() {
		defer echoServer.Close()
		for {
			conn, err := echoServer.Accept()
			if err != nil {
				log.Fatal(err)
			}
			go func() {
				defer conn.Close()
				buf := make([]byte, 1024)
				for {
					n, err := conn.Read(buf)
					if err != nil {
						return
					}
					fmt.Println("echo:", string(buf[:n]))
					conn.Write(buf[:n])
				}
			}()
		}
	}()

	// create a proxy server
	wsp := &WSproxy{
		URL:    "/echo",
		Target: "127.0.0.1:12340",
	}

	go func() {
		wsp.Serve("127.0.0.1:12341")
	}()

	// wait for the proxy server to start
	time.Sleep(100 * time.Millisecond)

	{
		// create a client
		r, err := http.Get("http://127.0.0.1:12341/")
		assert.Nil(t, err)
		assert.Equal(t, http.StatusNotFound, r.StatusCode)
	}
	{
		r, err := http.Get("http://127.0.0.1:12341/echo")
		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, r.StatusCode)
	}
	{
		// websocket client
		conn, _, err := websocket.DefaultDialer.Dial("ws://127.0.0.1:12341/echo", nil)
		assert.Nil(t, err)
		defer conn.Close()
		conn.WriteMessage(websocket.BinaryMessage, []byte("hello"))
		_, message, err := conn.ReadMessage()
		assert.Nil(t, err)
		assert.Equal(t, "hello", string(message))
	}
}
