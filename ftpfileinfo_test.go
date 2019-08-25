package graval

import (
	. "github.com/smartystreets/goconvey/convey"
	"os"
	"testing"
	"time"
)

func TestNewDirInfo(t *testing.T) {
	modTime := time.Unix(1566738000, 0) // 2019-08-25 13:00:00 UTC
	dirInfo := NewDirItem("dir", modTime)
	Convey("New Directory Info", t, func() {
		Convey("Will display the correct Mode", func() {
			So(dirInfo.Mode(), ShouldEqual, os.ModeDir|0666)
		})

		Convey("Will display the correct Name", func() {
			So(dirInfo.Name(), ShouldEqual, "dir")
		})

		Convey("Will display a size of 0 bytes", func() {
			So(dirInfo.Size(), ShouldEqual, 0)
		})

		Convey("Will display modified date as current time", func() {
			So(dirInfo.ModTime(), ShouldEqual, modTime)
		})

		Convey("Will return nil for Sys", func() {
			So(dirInfo.Sys(), ShouldBeNil)
		})
	})
}

func TestNewFileInfo(t *testing.T) {
	modTime := time.Unix(1566738000, 0) // 2019-08-25 13:00:00 UTC
	dirInfo := NewFileItem("test.txt", int64(99), modTime)
	Convey("New File Info", t, func() {
		Convey("Will display the correct Mode", func() {
			So(dirInfo.Mode(), ShouldEqual, 0666)
		})

		Convey("Will display the correct Name", func() {
			So(dirInfo.Name(), ShouldEqual, "test.txt")
		})

		Convey("Will display a size of 0 bytes", func() {
			So(dirInfo.Size(), ShouldEqual, 99)
		})

		Convey("Will display modified date as current time", func() {
			So(dirInfo.ModTime(), ShouldEqual, modTime)
		})

		Convey("Will return nil for Sys", func() {
			So(dirInfo.Sys(), ShouldBeNil)
		})
	})
}
