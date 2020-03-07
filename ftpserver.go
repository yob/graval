// An experimental FTP server framework. By providing a simple driver class that
// responds to a handful of methods you can have a complete FTP server.
//
// Some sample use cases include persisting data to an Amazon S3 bucket, a
// relational database, redis or memory.
//
// There is a sample in-memory driver available - see the documentation for the
// graval-mem package or the graval READEME for more details.
package graval

import (
	"net"
	"strconv"
	"strings"
	"time"
)

// serverOpts contains parameters for graval.NewFTPServer()
type FTPServerOpts struct {
	// Server name will be used for welcome message
	ServerName string

	// The factory that will be used to create a new FTPDriver instance for
	// each client connection. This is a mandatory option.
	Factory FTPDriverFactory

	// The hostname that the FTP server should listen on. Optional, defaults to
	// "::", which means all hostnames on ipv4 and ipv6.
	Hostname string

	// The port that the FTP should listen on. Optional, defaults to 3000. In
	// a production environment you will probably want to change this to 21.
	Port int

	// The lower bound of port numbers that can be used for passive-mode data sockets
	// Defaults to 0, which allows the server to pick any free port
	PasvMinPort int

	// The upper bound of port numbers that can be used for passive-mode data sockets
	// Defaults to 0, which allows the server to pick any free port
	PasvMaxPort int

	// Use this option to override the IP address that will be advertised in response to the
	// PASV command. Most setups can ignore this, but it can be helpful in situations where
	// the FTP server is behind a NAT gateway or load balancer and the public IP used by
	// clients is different to the IP the server is directly listening on
	PasvAdvertisedIp string
}

// FTPServer is the root of your FTP application. You should instantiate one
// of these and call ListenAndServe() to start accepting client connections.
//
// Always use the NewFTPServer() method to create a new FTPServer.
type FTPServer struct {
	serverName       string
	listenTo         string
	driverFactory    FTPDriverFactory
	logger           *ftpLogger
	pasvMinPort      int
	pasvMaxPort      int
	pasvAdvertisedIp string
	closeChan        chan struct{}
}

// serverOptsWithDefaults copies an FTPServerOpts struct into a new struct,
// then adds any default values that are missing and returns the new data.
func serverOptsWithDefaults(opts *FTPServerOpts) *FTPServerOpts {
	var newOpts FTPServerOpts

	if opts == nil {
		opts = &FTPServerOpts{}
	}

	if opts.ServerName == "" {
		newOpts.ServerName = "Go FTP Server"
	} else {
		newOpts.ServerName = opts.ServerName
	}

	if opts.Hostname == "" {
		newOpts.Hostname = "::"
	} else {
		newOpts.Hostname = opts.Hostname
	}

	if opts.Port == 0 {
		newOpts.Port = 3000
	} else {
		newOpts.Port = opts.Port
	}

	newOpts.PasvMinPort = opts.PasvMinPort
	newOpts.PasvMaxPort = opts.PasvMaxPort
	newOpts.PasvAdvertisedIp = opts.PasvAdvertisedIp
	newOpts.Factory = opts.Factory

	return &newOpts
}

// NewFTPServer initialises a new FTP server. Configuration options are provided
// via an instance of FTPServerOpts. Calling this function in your code will
// probably look something like this:
//
//     factory := &MyDriverFactory{}
//     server  := graval.NewFTPServer(&graval.FTPServerOpts{ Factory: factory })
//
// or:
//
//     factory := &MyDriverFactory{}
//     opts    := &graval.FTPServerOpts{
//       Factory: factory,
//       Port: 2000,
//       Hostname: "127.0.0.1",
//     }
//     server  := graval.NewFTPServer(opts)
//
func NewFTPServer(opts *FTPServerOpts) *FTPServer {
	opts = serverOptsWithDefaults(opts)
	s := new(FTPServer)
	s.listenTo = buildTcpString(opts.Hostname, opts.Port)
	s.serverName = opts.ServerName
	s.driverFactory = opts.Factory
	s.logger = newFtpLogger("")
	s.pasvMinPort = opts.PasvMinPort
	s.pasvMaxPort = opts.PasvMaxPort
	s.pasvAdvertisedIp = opts.PasvAdvertisedIp
	s.closeChan = make(chan struct{})
	return s
}

// ListenAndServe asks a new FTPServer to begin accepting client connections. It
// accepts no arguments - all configuration is provided via the NewFTPServer
// function.
//
// If the server fails to start for any reason, an error will be returned. Common
// errors are trying to bind to a privileged port or something else is already
// listening on the same port.
//
func (ftpServer *FTPServer) ListenAndServe() error {
	laddr, err := net.ResolveTCPAddr("tcp", ftpServer.listenTo)
	if err != nil {
		return err
	}
	listener, err := net.ListenTCP("tcp", laddr)
	if err != nil {
		return err
	}
	ftpServer.logger.Printf("listening on %s", listener.Addr().String())

	for {
		select {
		case <-ftpServer.closeChan:
			listener.Close()
			return nil
		default:
			listener.SetDeadline(time.Now().Add(2 * time.Second))
			tcpConn, err := listener.AcceptTCP()
			if strings.HasSuffix(err.Error(), "i/o timeout") {
				// deadline reached, no big deal
				// NOTE: This error is passed from the internal/poll/ErrTimeout but that
				// package is not legal to include, hence the string match. :(
				continue
			} else if err != nil {
				ftpServer.logger.Printf("listening error: %+v", err)
				return err
			}

			driver, err := ftpServer.driverFactory.NewDriver()
			if err != nil {
				ftpServer.logger.Print("Error creating driver, aborting client connection")
			} else {
				ftpConn := newftpConn(tcpConn, driver, ftpServer.serverName, ftpServer.pasvMinPort, ftpServer.pasvMaxPort, ftpServer.pasvAdvertisedIp)
				go ftpConn.Serve()
			}

		}
	}
	return nil
}

func (ftpServer *FTPServer) Close() {
	select {
	case <-ftpServer.closeChan:
	// already closed
	default:
		close(ftpServer.closeChan)
	}
}

func buildTcpString(hostname string, port int) (result string) {
	if strings.Contains(hostname, ":") {
		// ipv6
		if port == 0 {
			result = "[" + hostname + "]"
		} else {
			result = "[" + hostname + "]:" + strconv.Itoa(port)
		}
	} else {
		// ipv4
		if port == 0 {
			result = hostname
		} else {
			result = hostname + ":" + strconv.Itoa(port)
		}
	}
	return
}
