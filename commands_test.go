package graval

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestStringMapsToCorrectCommands(t *testing.T) {
	Convey("Command map calls correct objects", t, func() {
		So(commands["ALLO"], ShouldHaveSameTypeAs, commandAllo{})
		So(commands["CDUP"], ShouldHaveSameTypeAs, commandCdup{})
		So(commands["CWD"], ShouldHaveSameTypeAs, commandCwd{})
		So(commands["DELE"], ShouldHaveSameTypeAs, commandDele{})
		So(commands["EPRT"], ShouldHaveSameTypeAs, commandEprt{})
		So(commands["EPSV"], ShouldHaveSameTypeAs, commandEpsv{})
		So(commands["LIST"], ShouldHaveSameTypeAs, commandList{})
		So(commands["NLST"], ShouldHaveSameTypeAs, commandNlst{})
		So(commands["MDTM"], ShouldHaveSameTypeAs, commandMdtm{})
		So(commands["MKD"], ShouldHaveSameTypeAs, commandMkd{})
		So(commands["MODE"], ShouldHaveSameTypeAs, commandMode{})
		So(commands["NOOP"], ShouldHaveSameTypeAs, commandNoop{})
		So(commands["PASS"], ShouldHaveSameTypeAs, commandPass{})
		So(commands["PASV"], ShouldHaveSameTypeAs, commandPasv{})
		So(commands["PORT"], ShouldHaveSameTypeAs, commandPort{})
		So(commands["PWD"], ShouldHaveSameTypeAs, commandPwd{})
		So(commands["QUIT"], ShouldHaveSameTypeAs, commandQuit{})
		So(commands["RETR"], ShouldHaveSameTypeAs, commandRetr{})
		So(commands["RNFR"], ShouldHaveSameTypeAs, commandRnfr{})
		So(commands["RNTO"], ShouldHaveSameTypeAs, commandRnto{})
		So(commands["RMD"], ShouldHaveSameTypeAs, commandRmd{})
		So(commands["SIZE"], ShouldHaveSameTypeAs, commandSize{})
		So(commands["STOR"], ShouldHaveSameTypeAs, commandStor{})
		So(commands["STRU"], ShouldHaveSameTypeAs, commandStru{})
		So(commands["SYST"], ShouldHaveSameTypeAs, commandSyst{})
		So(commands["TYPE"], ShouldHaveSameTypeAs, commandType{})
		So(commands["USER"], ShouldHaveSameTypeAs, commandUser{})
		So(commands["XCUP"], ShouldHaveSameTypeAs, commandCdup{})
		So(commands["XCWD"], ShouldHaveSameTypeAs, commandCwd{})
		So(commands["XPWD"], ShouldHaveSameTypeAs, commandPwd{})
		So(commands["XRMD"], ShouldHaveSameTypeAs, commandRmd{})
	})
}
