package graval

import (
	"errors"
	"fmt"
	"net"
	"time"
)

// A data socket is used to send non-control data between the client and
// server.
type ftpDataSocket interface {
	Host() string

	Port() int

	// the standard io.Reader interface
	Read(p []byte) (n int, err error)

	// the standard io.Writer interface
	Write(p []byte) (n int, err error)

	// the standard io.Closer interface
	Close() error
}

type ftpActiveSocket struct {
	conn *net.TCPConn
	host string
	port int
	logger *ftpLogger
}

func newActiveSocket(host string, port int, logger *ftpLogger) (ftpDataSocket, error) {
	connectTo := buildTcpString(host, port)
	logger.Print("Opening active data connection to " + connectTo)
	raddr, err := net.ResolveTCPAddr("tcp", connectTo)
	if err != nil {
		logger.Print(err)
		return nil, err
	}
	tcpConn, err := net.DialTCP("tcp", nil, raddr)
	if err != nil {
		logger.Print(err)
		return nil, err
	}
	socket := new(ftpActiveSocket)
	socket.conn = tcpConn
	socket.host = host
	socket.port = port
	socket.logger = logger
	return socket, nil
}

func (socket *ftpActiveSocket) Host() string {
	return socket.host
}

func (socket *ftpActiveSocket) Port() int {
	return socket.port
}

func (socket *ftpActiveSocket) Read(p []byte) (n int, err error) {
	return socket.conn.Read(p)
}

func (socket *ftpActiveSocket) Write(p []byte) (n int, err error) {
	return socket.conn.Write(p)
}

func (socket *ftpActiveSocket) Close() error {
	return socket.conn.Close()
}


type ftpPassiveSocket struct {
	conn     *net.TCPConn
	port     int
	ingress  chan []byte
	egress   chan []byte
	logger   *ftpLogger
}

func newPassiveSocket(logger *ftpLogger, minPort int, maxPort int) (ftpDataSocket, error) {
	socket := new(ftpPassiveSocket)
	socket.ingress = make(chan []byte)
	socket.egress = make(chan []byte)
	socket.logger = logger
	go socket.ListenAndServe(minPort, maxPort)
	for {
		if socket.Port() > 0 {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
	return socket, nil
}

func (socket *ftpPassiveSocket) Host() string {
	return "127.0.0.1"
}

func (socket *ftpPassiveSocket) Port() int {
	return socket.port
}

func (socket *ftpPassiveSocket) Read(p []byte) (n int, err error) {
	if socket.waitForOpenSocket() == false {
		return 0, errors.New("data socket unavailable")
	}
	return socket.conn.Read(p)
}

func (socket *ftpPassiveSocket) Write(p []byte) (n int, err error) {
	if socket.waitForOpenSocket() == false {
		return 0, errors.New("data socket unavailable")
	}
	return socket.conn.Write(p)
}

func (socket *ftpPassiveSocket) Close() error {
	socket.logger.Print("closing passive data socket")
	return socket.conn.Close()
}

func (socket *ftpPassiveSocket) ListenAndServe(minPort int, maxPort int) {
	listener, err := socket.netListenerInRange(minPort, maxPort)
	if err != nil {
		socket.logger.Print(err)
		return
	}
	add   := listener.Addr().(*net.TCPAddr)
	socket.port = add.Port
	tcpConn, err := listener.AcceptTCP()
	if err != nil {
		socket.logger.Print(err)
		return
	}
	socket.conn = tcpConn
}

func (socket *ftpPassiveSocket) waitForOpenSocket() bool {
	retries := 0
	for {
		if socket.conn != nil {
			break
		}
		if retries > 3 {
			return false
		}
		socket.logger.Print("sleeping, socket isn't open")
		sleepMs := time.Duration(500 * (retries+1))
		time.Sleep(sleepMs *time.Millisecond)
		retries += 1
	}
	return true
}

func (socket *ftpPassiveSocket) netListenerInRange(min, max int) (*net.TCPListener, error) {
	socket.logger.Printf("looking for new socket in range %d-%d", min, max)
	for port := min; port <= max; port++ {
		l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
		if err == nil {
			return l.(*net.TCPListener), nil
		}
	}

	return nil, errors.New("Unable to find available port to listen on")
}
