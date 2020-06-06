package graval

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

func TestClose(t *testing.T) {

	goneChan := make(chan struct{})

	Convey("Setting up a minimal server, it will end if Close() is called ", t, func() {
		opts := &FTPServerOpts{
			ServerName:  "blah blah blah",
			PasvMinPort: 60200,
			PasvMaxPort: 60300,
		}
		ftpServer := NewFTPServer(opts)
		go func() {
			defer close(goneChan)
			err := ftpServer.ListenAndServe()
			if err != nil {
				panic(err)
			}
		}()
		time.Sleep(1 * time.Second)
		So(ftpServer.Close, ShouldNotPanic)
		<-goneChan

	})
}

func TestCloseLater(t *testing.T) {
	t.Skip("Waits a while, which you don't need to")
	goneChan := make(chan struct{})

	Convey("Setting up a minimal server, it will end if Close() is called ", t, func() {
		opts := &FTPServerOpts{
			ServerName:  "blah blah blah",
			PasvMinPort: 60200,
			PasvMaxPort: 60300,
		}
		ftpServer := NewFTPServer(opts)
		go func() {
			defer close(goneChan)
			err := ftpServer.ListenAndServe()
			if err != nil {
				panic(err)
			}
		}()
		time.Sleep(10 * time.Second)
		So(ftpServer.Close, ShouldNotPanic)
		<-goneChan

	})
}

func TestCloseHammer(t *testing.T) {

	goneChan := make(chan struct{})

	Convey("Setting up a minimal server, it will end if Close() is called a bunch of times", t, func() {
		opts := &FTPServerOpts{
			ServerName:  "blah blah blah",
			PasvMinPort: 60200,
			PasvMaxPort: 60300,
		}
		ftpServer := NewFTPServer(opts)
		go func() {
			defer close(goneChan)
			err := ftpServer.ListenAndServe()
			if err != nil {
				panic(err)
			}
		}()
		time.Sleep(1 * time.Second)
		So(ftpServer.Close, ShouldNotPanic)
		So(ftpServer.Close, ShouldNotPanic)
		So(ftpServer.Close, ShouldNotPanic)
		So(ftpServer.Close, ShouldNotPanic)
		So(ftpServer.Close, ShouldNotPanic)
		So(ftpServer.Close, ShouldNotPanic)
		So(ftpServer.Close, ShouldNotPanic)
		So(ftpServer.Close, ShouldNotPanic)

		<-goneChan

	})
}
