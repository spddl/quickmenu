// Copyright 2011 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

package walk

import (
	"github.com/lxn/win"
)

func (ni *NotifyIcon) OpenContextMenu(hwnd win.HWND) {
	var p win.POINT
	if !win.GetCursorPos(&p) {
		lastError("GetCursorPos")
	}
	ni.applyDPI()

	actionID := uint16(win.TrackPopupMenuEx(
		ni.contextMenu.hMenu,
		win.TPM_NOANIMATION|win.TPM_RETURNCMD,
		p.X,
		p.Y,
		hwnd,
		nil))

	if actionID != 0 {
		if action, ok := actionsById[actionID]; ok {
			action.raiseTriggered()
		}
	}
}
