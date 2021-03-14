// https://stackoverflow.com/questions/38191972/call-shgetimagelist-in-go
package main

import (
	"syscall"
	"unsafe"

	"github.com/lxn/win"
)

func hIconForFilePath(filePath string) win.HICON {
	fPptr, _ := syscall.UTF16PtrFromString(filePath)
	var shfi win.SHFILEINFO
	hIml := win.HIMAGELIST(win.SHGetFileInfo( // https://docs.microsoft.com/en-us/windows/win32/api/shellapi/nf-shellapi-shgetfileinfow
		fPptr,
		0,
		&shfi,
		uint32(unsafe.Sizeof(shfi)),
		win.SHGFI_ICON|win.SHGFI_SMALLICON))
	if hIml != 0 {
		return shfi.HIcon
	}
	return 0
}
