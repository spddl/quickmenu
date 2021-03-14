package main

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/karrick/godirwalk"
	"github.com/lxn/walk"

	"github.com/lxn/win"
)

var dpi int
var dir string

func main() {
	var err error
	if len(os.Args) == 1 {
		dir, err = os.Getwd()
		if err != nil {
			panic(err)
		}
	} else {
		dir = os.Args[1]
	}

	// We need either a walk.MainWindow or a walk.Dialog for their message loop.
	// We will not make it visible in this example, though.
	mw, err := walk.NewMainWindow()
	if err != nil {
		panic(err)
	}

	// Create the notify icon and make sure we clean it up on exit.
	ni, err := walk.NewNotifyIcon(mw)
	if err != nil {
		panic(err)
	}
	defer ni.Dispose()

	files, err := osReadDir(dir)
	if err != nil {
		panic(err)
	}

	dpi = screenDPI()
	ActionList := ni.ContextMenu().Actions()
	for _, file := range files {
		createMenuEntry(ActionList, file)
	}

	err = ActionList.Add(walk.NewSeparatorAction())
	if err != nil {
		panic(err)
	}

	ExitMenuItem := walk.NewAction()
	err = ExitMenuItem.SetText("Exit")
	if err != nil {
		panic(err)
	}
	ExitMenuItem.Triggered().Attach(func() {
		walk.App().Exit(0)
	})
	err = ActionList.Add(ExitMenuItem)
	if err != nil {
		panic(err)
	}

	// The notify icon is hidden initially, so we have to make it visible.
	if err := ni.SetVisible(true); err != nil {
		panic(err)
	}

	ni.OpenContextMenu(mw.Handle())
}

func createMenuEntry(al *walk.ActionList, file string) {
	menuItem := walk.NewAction()

	hI := hIconForFilePath(filepath.Join(dir, file))
	if hI != 0 {
		ic, err := walk.NewIconFromHICONForDPI(hI, dpi)
		if err == nil {
			menuItem.SetImage(ic)
		} else {
			panic(err)
		}
	}

	if err := menuItem.SetText(file); err != nil {
		panic(err)
	}

	menuItem.Triggered().Attach(func() {
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

		exec.Command(prog, args...).Start()
	})
	if err := al.Add(menuItem); err != nil {
		panic(err)
	}
}

func osReadDir(dirname string) ([]string, error) {
	var files []string
	err := godirwalk.Walk(dirname, &godirwalk.Options{
		Callback: func(osPathname string, de *godirwalk.Dirent) error {
			if de.IsDir() {
				rel, err := filepath.Rel(dirname, osPathname)
				if err != nil {
					panic(err)
				}
				if rel != "." {
					return godirwalk.SkipThis
				}
				return nil
			}
			files = append(files, de.Name())
			return nil
		},
		Unsorted: false,
	})
	if err != nil {
		panic(err)
	}
	return files, nil
}

func screenDPI() int {
	hDC := win.GetDC(0)
	defer win.ReleaseDC(0, hDC)
	return int(win.GetDeviceCaps(hDC, win.LOGPIXELSY))
}
