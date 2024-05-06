# A Websockify with 100 lines of Golang code

**websockify-go** is a pure Go implementation of [novnc/websockify](https://github.com/novnc/websockify) TCP to WebSocket proxy with improved connection handling.

By [ruzhila.cn](http://ruzhila.cn/?from=github_websockify), a campus for learning backend development through practice.

This is a tutorial code demonstrating how to use Golang write network proxy. Pull requests are welcome. 👏

## Features
- 👏 Simple and easy to use
- 🐶 Pure Go implementation
- 🔐 Support SSL
- ⌨️ Support custom URL path
- 🚪 Support custom target address

## Build from source
```shell
$ git clone https://github.com/ruzhila/websockify-go.git
$ cd websockify-go/cmd
$ go build -o websockify main.go
```
## Install from source
```shell
$ go get -u github.com/ruzhila/websockify-go
$ websockifygo -h
```

## Usage
```shell
$ ./websockify -h
Usage of ./websockify:
  -addr string
        address to listen on (default ":8080")
  -cert string
        SSL cert.pem
  -key string
        SSL key.pem
  -target string
        target address (default "localhost:5900")
  -url string
        url path to proxy, e.g. /vnc

$ ./websockify -addr :8080 -target localhost:5900 -url /vnc
```
### Connect to the proxy, via browser
```javascript
const ws = new WebSocket('ws://localhost:8080/vnc');
ws.onopen = () => console.log('connected');
ws.onmessage = (e) => console.log('message', e.data);
ws.onclose = () => console.log('closed');
```

## Build with Docker
```shell
$ docker build -t websockify-go .
$ docker run -p 8080:8080 websockify-go -target localhost:5900 -url /vnc
```