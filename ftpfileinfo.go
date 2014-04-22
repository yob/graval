package graval

import (
	"os"
	"time"
)

type ftpFileInfo struct {
	name    string
	bytes   int64
	mode    os.FileMode
	modtime time.Time
}

func (info *ftpFileInfo) Name() string {
	return info.name
}

func (info *ftpFileInfo) Size() int64 {
	return info.bytes
}

func (info *ftpFileInfo) Mode() os.FileMode {
	return info.mode
}

func (info *ftpFileInfo) ModTime() time.Time {
	return info.modtime
}

func (info *ftpFileInfo) IsDir() bool {
	return (info.mode | os.ModeDir) == os.ModeDir
}

func (info *ftpFileInfo) Sys() interface{} {
	return nil
}

// NewDirItem creates a new os.FileInfo that represents a single diretory. Use
// this function to build the response to DirContents() in your FTPDriver
// implementation.
func NewDirItem(name string) os.FileInfo {
	d := new(ftpFileInfo)
	d.name = name
	d.bytes = int64(0)
	d.mode = os.ModeDir | 0666
	d.modtime = time.Now().UTC()
	return d
}

// NewFileItem creates a new os.FileInfo that represents a single file. Use
// this function to build the response to DirContents() in your FTPDriver
// implementation.
func NewFileItem(name string, bytes int64, modtime time.Time) os.FileInfo {
	f := new(ftpFileInfo)
	f.name = name
	f.bytes = int64(bytes)
	f.mode = 0666
	f.modtime = modtime
	return f
}
