// https://stackoverflow.com/questions/38191972/call-shgetimagelist-in-go
package main

import (
	"log"
	"path/filepath"
	"sync"
	"syscall"
	"unsafe"

	"github.com/lxn/walk"
	"github.com/lxn/win"
)

func hIconForFilePath(filePath string) (string, win.HICON) {
	fPptr, err := syscall.UTF16PtrFromString(filePath)
	ErrCheck(err)
	var shfi win.SHFILEINFO
	hIml := win.HIMAGELIST(win.SHGetFileInfo( // https://docs.microsoft.com/en-us/windows/win32/api/shellapi/nf-shellapi-shgetfileinfow
		fPptr,
		0x80, // FILE_ATTRIBUTE_NORMAL
		&shfi,
		uint32(unsafe.Sizeof(shfi)),
		win.SHGFI_USEFILEATTRIBUTES|win.SHGFI_ICON|win.SHGFI_SMALLICON|win.SHGFI_DISPLAYNAME))
	if hIml != 0 {
		return syscall.UTF16ToString(shfi.SzDisplayName[:]), shfi.HIcon
	}
	return "", 0
}

func (io *IconsObject) getIcons(wg *sync.WaitGroup, file string) {
	defer wg.Done()
	io.m.Lock() // that seems necessary
	displayname, hI := hIconForFilePath(filepath.Join(dir, file))
	io.m.Unlock()
	if hI != 0 {
		ic, err := walk.NewIconFromHICONForDPI(hI, dpi)
		if err != nil {
			log.Println("getIcons Err", err)
			return
		}
		io.m.Lock()
		io.ic[file] = FileInfo{
			DisplayName: displayname,
			Icon:        ic,
		}
		io.m.Unlock()
	}
}
