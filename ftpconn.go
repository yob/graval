package graval

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net"
	"path/filepath"
	"strings"
	"time"
)

type ftpConn struct {
	conn          *net.TCPConn
	controlReader *bufio.Reader
	controlWriter *bufio.Writer
	dataConn      ftpDataSocket
	driver        FTPDriver
	logger        *ftpLogger
	serverName    string
	sessionId     string
	namePrefix    string
	reqUser       string
	user          string
	renameFrom    string
	minDataPort   int
	maxDataPort   int
}

// NewftpConn constructs a new object that will handle the FTP protocol over
// an active net.TCPConn. The TCP connection should already be open before
// it is handed to this functions. driver is an instance of FTPDriver that
// will handle all auth and persistence details.
func newftpConn(tcpConn *net.TCPConn, driver FTPDriver, serverName string, minPort int, maxPort int) *ftpConn {
	c := new(ftpConn)
	c.namePrefix = "/"
	c.conn = tcpConn
	c.controlReader = bufio.NewReader(tcpConn)
	c.controlWriter = bufio.NewWriter(tcpConn)
	c.driver = driver
	c.sessionId = newSessionId()
	c.logger = newFtpLogger(c.sessionId)
	c.serverName = serverName
	c.minDataPort = minPort
	c.maxDataPort = maxPort
	return c
}

// returns a random 20 char string that can be used as a unique session ID
func newSessionId() string {
	hash := sha256.New()
	_, err := io.CopyN(hash, rand.Reader, 50)
	if err != nil {
		return "????????????????????"
	}
	md := hash.Sum(nil)
	mdStr := hex.EncodeToString(md)
	return mdStr[0:20]
}

// Serve starts an endless loop that reads FTP commands from the client and
// responds appropriately. terminated is a channel that will receive a true
// message when the connection closes. This loop will be running inside a
// goroutine, so use this channel to be notified when the connection can be
// cleaned up.
func (ftpConn *ftpConn) Serve() {
	defer func() {
		if r := recover(); r != nil {
			ftpConn.logger.Printf("Recovered in ftpConn Serve: %s", r)
		}

		ftpConn.Close()
	}()

	ftpConn.logger.Printf("Connection Established (local: %s, remote: %s)", ftpConn.localIP(), ftpConn.remoteIP())
	// send welcome
	ftpConn.writeMessage(220, ftpConn.serverName)
	// read commands
	for {
		line, err := ftpConn.controlReader.ReadString('\n')
		if err != nil {
			break
		}
		ftpConn.receiveLine(line)
	}
	ftpConn.logger.Print("Connection Terminated")
}

// Close will manually close this connection, even if the client isn't ready.
func (ftpConn *ftpConn) Close() {
	ftpConn.conn.Close()
	if ftpConn.dataConn != nil {
		ftpConn.dataConn.Close()
	}
}

// receiveLine accepts a single line FTP command and co-ordinates an
// appropriate response.
func (ftpConn *ftpConn) receiveLine(line string) {
	command, param := ftpConn.parseLine(line)
	ftpConn.logger.PrintCommand(command, param)
	cmdObj := commands[command]
	if cmdObj == nil {
		ftpConn.writeMessage(500, "Command not found")
		return
	}
	if cmdObj.RequireParam() && param == "" {
		ftpConn.writeMessage(553, "action aborted, required param missing")
	} else if cmdObj.RequireAuth() && ftpConn.user == "" {
		ftpConn.writeMessage(530, "not logged in")
	} else {
		cmdObj.Execute(ftpConn, param)
	}
}

func (ftpConn *ftpConn) parseLine(line string) (string, string) {
	params := strings.SplitN(strings.Trim(line, "\r\n"), " ", 2)
	if len(params) == 1 {
		return params[0], ""
	}
	return params[0], strings.TrimSpace(params[1])
}

// writeMessage will send a standard FTP response back to the client.
func (ftpConn *ftpConn) writeMessage(code int, message string) (wrote int, err error) {
	ftpConn.logger.PrintResponse(code, message)
	line := fmt.Sprintf("%d %s\r\n", code, message)
	wrote, err = ftpConn.controlWriter.WriteString(line)
	ftpConn.controlWriter.Flush()
	return
}

// buildPath takes a client supplied path or filename and generates a safe
// absolute path within their account sandbox.
//
//    buildpath("/")
//    => "/"
//    buildpath("one.txt")
//    => "/one.txt"
//    buildpath("/files/two.txt")
//    => "/files/two.txt"
//    buildpath("files/two.txt")
//    => "files/two.txt"
//    buildpath("/../../../../etc/passwd")
//    => "/etc/passwd"
//
// The driver implementation is responsible for deciding how to treat this path.
// Obviously they MUST NOT just read the path off disk. The probably want to
// prefix the path with something to scope the users access to a sandbox.
func (ftpConn *ftpConn) buildPath(filename string) (fullPath string) {
	if len(filename) > 0 && filename[0:1] == "/" {
		fullPath = filepath.Clean(filename)
	} else if len(filename) > 0 && filename != "-a" {
		fullPath = filepath.Clean(ftpConn.namePrefix + "/" + filename)
	} else {
		fullPath = filepath.Clean(ftpConn.namePrefix)
	}
	fullPath = strings.Replace(fullPath, "//", "/", -1)
	return
}

// the server IP that is being used for this connection. May be the same for all connections,
// or may vary if the server is listening on 0.0.0.0
func (ftpConn *ftpConn) localIP() string {
	lAddr := ftpConn.conn.LocalAddr().(*net.TCPAddr)
	return lAddr.IP.String()
}

// the client IP address
func (ftpConn *ftpConn) remoteIP() string {
	rAddr := ftpConn.conn.RemoteAddr().(*net.TCPAddr)
	return rAddr.IP.String()
}

// sendOutofbandData will copy data from reader to the client via the currently
// open data socket. Assumes the socket is open and ready to be used.
func (ftpConn *ftpConn) sendOutofbandReader(reader io.Reader) {
	defer ftpConn.dataConn.Close()

	_, err := io.Copy(ftpConn.dataConn, reader)

	if err != nil {
		ftpConn.logger.Printf("sendOutofbandReader copy error %s", err)
		ftpConn.writeMessage(550, "Action not taken")
		return
	}

	ftpConn.writeMessage(226, "Transfer complete.")

	// Chrome dies on localhost if we close connection to soon
	time.Sleep(10 * time.Millisecond)
}

// sendOutofbandData will send a string to the client via the currently open
// data socket. Assumes the socket is open and ready to be used.
func (ftpConn *ftpConn) sendOutofbandData(data string) {
	ftpConn.sendOutofbandReader(bytes.NewReader([]byte(data)))
}

func (ftpConn *ftpConn) newPassiveSocket() (socket *ftpPassiveSocket, err error) {
	if ftpConn.dataConn != nil {
		ftpConn.dataConn.Close()
		ftpConn.dataConn = nil
	}

	socket, err = newPassiveSocket(ftpConn.localIP(), ftpConn.minDataPort, ftpConn.maxDataPort, ftpConn.logger)

	if err == nil {
		ftpConn.dataConn = socket
	}

	return
}

func (ftpConn *ftpConn) newActiveSocket(host string, port int) (socket *ftpActiveSocket, err error) {
	if ftpConn.dataConn != nil {
		ftpConn.dataConn.Close()
		ftpConn.dataConn = nil
	}

	socket, err = newActiveSocket(host, port, ftpConn.logger)

	if err == nil {
		ftpConn.dataConn = socket
	}

	return
}
