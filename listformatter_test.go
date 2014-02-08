package graval

import (
	. "github.com/smartystreets/goconvey/convey"
	"os"
	"testing"
	"time"
)

type TestFileInfo struct{}

func (t *TestFileInfo) Name() string {
	return "file1.txt"
}

func (t *TestFileInfo) Size() int64 {
	return 99
}

func (t *TestFileInfo) Mode() os.FileMode {
	return os.ModeSymlink
}

func (t *TestFileInfo) IsDir() bool {
	return false
}

func (t *TestFileInfo) ModTime() time.Time {
	return time.Now()
}

func (t *TestFileInfo) Sys() interface{} {
	return nil
}

var files []os.FileInfo = []os.FileInfo{
	&TestFileInfo{}, &TestFileInfo{},
}

func TestShortFormat(t *testing.T) {
	formatter := newListFormatter(files)
	Convey("The Short listing format", t, func() {
		Convey("Will display correctly", func() {
			So(formatter.Short(), ShouldEqual, "file1.txt\r\nfile1.txt\r\n\r\n")
		})
	})
}
