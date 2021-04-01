package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"sync"

	"github.com/karrick/godirwalk"
	"github.com/lxn/walk"

	"github.com/lxn/win"
)

var dpi int
var dir string

type FileInfo struct {
	Icon        *walk.Icon
	DisplayName string
}

type IconsObject struct {
	m  sync.RWMutex
	ic map[string]FileInfo
}

func main() {
	var err error
	if len(os.Args) == 1 {
		dir, err = os.Getwd()
		ErrCheck(err)
	} else {
		dir = os.Args[1]
	}

	// We need either a walk.MainWindow or a walk.Dialog for their message loop.
	// We will not make it visible in this example, though.
	mw, err := walk.NewMainWindow()
	ErrCheck(err)

	// Create the notify icon and make sure we clean it up on exit.
	ni, err := walk.NewNotifyIcon(mw)
	ErrCheck(err)
	defer ni.Dispose()

	files, err := osReadDir(dir)
	ErrCheck(err)

	dpi = screenDPI()

	ActionList := ni.ContextMenu().Actions()

	io := IconsObject{
		ic: make(map[string]FileInfo, len(files)),
	}
	var wg sync.WaitGroup
	wg.Add(len(files))
	for _, file := range files {
		go io.getIcons(&wg, file)
	}
	wg.Wait()

	for _, file := range files {
		io.createMenuEntry(ActionList, file)
	}

	ErrCheck(ActionList.Add(walk.NewSeparatorAction()))

	ExitMenuItem := walk.NewAction()
	ErrCheck(ExitMenuItem.SetText("Exit"))
	ExitMenuItem.Triggered().Attach(func() {
		walk.App().Exit(0)
	})

	ErrCheck(ActionList.Add(ExitMenuItem))

	// The notify icon is hidden initially, so we have to make it visible.
	ErrCheck(ni.SetVisible(true))

	ni.OpenContextMenu(mw.Handle())
}

func ErrCheck(err error) {
	if err != nil {
		panic(err)
	}
}

func (io *IconsObject) createMenuEntry(al *walk.ActionList, file string) {
	menuItem := walk.NewAction()

	if val, ok := io.ic[file]; ok {
		ErrCheck(menuItem.SetImage(val.Icon))
		ErrCheck(menuItem.SetText(val.DisplayName))
	} else {
		ErrCheck(menuItem.SetText(file))
	}

	menuItem.Triggered().Attach(createCommandCall(file))
	ErrCheck(al.Add(menuItem))
}

func osReadDir(dirname string) ([]string, error) {
	var files []string
	ErrCheck(godirwalk.Walk(dirname, &godirwalk.Options{
		Callback: func(osPathname string, de *godirwalk.Dirent) error {
			if de.IsDir() {
				rel, err := filepath.Rel(dirname, osPathname)
				ErrCheck(err)
				if rel != "." {
					return godirwalk.SkipThis
				}
				return nil
			}
			files = append(files, de.Name())
			return nil
		},
		Unsorted: false,
	}))
	return files, nil
}

func screenDPI() int {
	hDC := win.GetDC(0)
	defer win.ReleaseDC(0, hDC)
	return int(win.GetDeviceCaps(hDC, win.LOGPIXELSY))
}

func createCommandCall(file string) func() {
	var prog string
	var args []string
	switch filepath.Ext(file) {
	case ".ps1": // preset
		prog = "PowerShell"
		args = []string{"-NoLogo", "-NoProfile", "-ExecutionPolicy Bypass", "-File", filepath.Join(dir, file)} // https://docs.microsoft.com/de-de/powershell/module/microsoft.powershell.core/about/about_powershell_exe?view=powershell-5.1

	default: // https://stackoverflow.com/a/12076082
		prog = "rundll32.exe"
		args = []string{"url.dll,FileProtocolHandler", filepath.Join(dir, file)}
	}

	return func() {
		exec.Command(prog, args...).Start()
	}
}
