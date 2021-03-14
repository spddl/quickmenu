module github.com/spddl/quickmenu

go 1.16

replace github.com/lxn/walk => ./update/github.com/lxn/walk

require (
	github.com/karrick/godirwalk v1.16.1
	github.com/lxn/walk v0.0.0-20210112085537-c389da54e794
	github.com/lxn/win v0.0.0-20210218163916-a377121e959e
	golang.org/x/sys v0.0.0-20210313202042-bd2e13477e9c // indirect
)
