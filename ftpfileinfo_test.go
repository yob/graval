package graval

import (
	. "github.com/smartystreets/goconvey/convey"
	"os"
	"testing"
	"time"
)

func TestNewDirInfo(t *testing.T) {
	dirInfo := NewDirItem("dir")
	Convey("New Directory Info", t, func() {
		Convey("Will display the correct Mode", func() {
			So(dirInfo.Mode(), ShouldEqual, os.ModeDir|666)
		})

		Convey("Will display the correct Name", func() {
			So(dirInfo.Name(), ShouldEqual, "dir")
		})

		Convey("Will display a size of 0 bytes", func() {
			So(dirInfo.Size(), ShouldEqual, 0)
		})

		Convey("Will display modified date as current time", func() {
			beforeTime := time.Now()
			So(dirInfo.ModTime(), ShouldHappenBetween, beforeTime, time.Now())
		})

		Convey("Will return nil for Sys", func() {
			So(dirInfo.Sys(), ShouldBeNil)
		})
	})
}

func TestNewFileInfo(t *testing.T) {
	dirInfo := NewFileItem("test.txt", 99)
	Convey("New File Info", t, func() {
		Convey("Will display the correct Mode", func() {
			So(dirInfo.Mode(), ShouldEqual, 666)
		})

		Convey("Will display the correct Name", func() {
			So(dirInfo.Name(), ShouldEqual, "test.txt")
		})

		Convey("Will display a size of 0 bytes", func() {
			So(dirInfo.Size(), ShouldEqual, 99)
		})

		Convey("Will display modified date as current time", func() {
			beforeTime := time.Now()
			So(dirInfo.ModTime(), ShouldHappenBetween, beforeTime, time.Now())
		})

		Convey("Will return nil for Sys", func() {
			So(dirInfo.Sys(), ShouldBeNil)
		})
	})
}
