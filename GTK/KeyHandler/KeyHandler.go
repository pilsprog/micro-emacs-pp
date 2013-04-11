package keyhandler

import (
	"C"
	"github.com/mattn/go-gtk/gdk"
	"micro-emacs-pp/editor"
	. "micro-emacs-pp/keyhandler"
)

var (
	CtrlXHandler KeyHandler = GuardHandler(
		KeyCtrlx,
		PauseHandler(KeyChoice(CtrlFHandler, CtrlSHandler)))
	// CtrlFhandler Opens a file if Ctrl+F was
	// pressed.
	CtrlFHandler KeyHandler = GuardHandler(
		KeyCtrlf,
		ActionHandler(func(e *editor.Editor) KeyHandler {
			e.CommandBuf.GrabFocus()
			e.CommandBuf.Clear()
			e.CommandBuf.Write([]byte("Find-file:"))
			return InputHandler(func(s string, e *editor.Editor) KeyHandler {
				e.OpenFile(s[10:])
				e.CommandBuf.Clear()
				e.CommandBuf.Write([]byte("File Opened!"))
				e.Buf.GrabFocus()
				return Root
			})
		}))

	//CtrlSHandler saves the current buffer if
	//ctrl+S was pressed.
	CtrlSHandler KeyHandler = GuardHandler(
		KeyCtrls,
		ActionHandler(func(e *editor.Editor) KeyHandler {
			e.CommandBuf.GrabFocus()
			e.CommandBuf.Clear()
			return InputHandler(func(s string, e *editor.Editor) KeyHandler {
				e.SaveFile(s)
				e.CommandBuf.Clear()
				e.CommandBuf.Write([]byte("File Saved!"))
				e.Buf.GrabFocus()
				return Root
			})
		}))

	KeyCtrle KeyPressEvent = GTKKeyPressEvent{gdk.KEY_e, gdk.CONTROL_MASK}
	KeyCtrlx KeyPressEvent = GTKKeyPressEvent{gdk.KEY_x, gdk.CONTROL_MASK}
	KeyCtrlf KeyPressEvent = GTKKeyPressEvent{gdk.KEY_f, gdk.CONTROL_MASK}
	KeyCtrls KeyPressEvent = GTKKeyPressEvent{gdk.KEY_s, gdk.CONTROL_MASK}
)

// KeyPressEvent represents a keypress consisting of the particular key
// (KeyVal) and possibly a modifier (0 if no modifier is given).
type GTKKeyPressEvent struct {
	KeyVal   int
	Modifier gdk.ModifierType
}

func (press GTKKeyPressEvent) GetKeyValue() int {
	return press.KeyVal
}

func (press GTKKeyPressEvent) GetModifier() Modifier {
	mod := press.Modifier
	switch {
	case mod&gdk.CONTROL_MASK != 0:
		return CTRL
	}
	return NONE
}

// Compare two KeyPressEvents
func (k1 GTKKeyPressEvent) Equals(k2 KeyPressEvent) bool {
	return k1.GetKeyValue() == k2.GetKeyValue() && k1.GetModifier() == k2.GetModifier()
}
