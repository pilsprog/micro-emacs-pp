// Package Keyhandler provides combinators for
// writing corded keyevents with calls to microemacspp/Editor. The combinators
// builds a tree of possible key combinations.
package KeyHandler

import (
	"C"
	"fmt"
	"github.com/mattn/go-gtk/gdk"
	"micro-emacs-pp/Editor"
)

// KeyPressEvent represents a keypress consisting of the particular key
// (KeyVal) and possibly a modifier (0 if no modifier is given).
type KeyPressEvent struct {
	KeyVal   int
	Modifier gdk.ModifierType
}

// Compare two KeyPressEvents
func (k1 KeyPressEvent) Equals(k2 KeyPressEvent) bool {
	return k1.KeyVal == k2.KeyVal && k1.Modifier == k2.Modifier
}

// KeyHandler is the interface that every node in the tree implements. A
// KeyHandler takes a KeyPressEvent (the key that was pressed) the Editor and
// returns a boolean and a new KeyHandler.  true indicates that the keypress
// was successfully applied, and the keyHandler is a pointer to the new
// position in the tree.
type KeyHandler interface {
	Handle(e KeyPressEvent, editor *Editor.Editor) (bool, KeyHandler)
}

// KeyChoice is a KeyHandler which represents a choice between several different
// keyHandlers. KeyChoice calls every KeyHandler in its list until one succeeds.
type KeyChoice struct {
	choices []KeyHandler
}

func (k *KeyChoice) Handle(e KeyPressEvent, editor *Editor.Editor) (bool, KeyHandler) {
	for _, choice := range k.choices {
		ok, handler := choice.Handle(e, editor)
		if ok {
			return ok, handler
		}
	}
	return false, root
}

var (
	root         KeyHandler
	CtrlXHandler KeyHandler = &GuardHandler{
		KeyCtrlx,
		&PauseHandler{&KeyChoice{[]KeyHandler{
			CtrlFHandler,
			CtrlSHandler}}}}
	// CtrlFhandler Opens a file if Ctrl+F was
	// pressed.
	CtrlFHandler KeyHandler = &GuardHandler{
		KeyCtrlf,
		&ActionHandler{func(e *Editor.Editor) KeyHandler {
			e.CommandBuf.GrabFocus()
			e.CommandBuf.Clear()
			e.CommandBuf.Write([]byte("Find-file:"))
			return &InputHandler{"", func(s string, e *Editor.Editor) KeyHandler {
				fmt.Println(s)
				e.OpenFile(s)
				e.CommandBuf.Clear()
				e.CommandBuf.Write([]byte("File Opened!"))
				e.Buf.GrabFocus()
				return root
			}}
		}}}

	//CtrlSHandler saves the current buffer if
	//Ctrl+S was pressed.
	CtrlSHandler KeyHandler = &GuardHandler{
		KeyCtrls,
		&ActionHandler{func(e *Editor.Editor) KeyHandler {
			e.CommandBuf.GrabFocus()
			e.CommandBuf.Clear()
			return &InputHandler{"", func(s string, e *Editor.Editor) KeyHandler {
				e.SaveFile(s)
				e.CommandBuf.Clear()
				e.CommandBuf.Write([]byte("File Saved!"))
				e.Buf.GrabFocus()
				return root
			}}
		}}}

	KeyReturn KeyPressEvent = KeyPressEvent{gdk.KEY_Return, 0}
	KeyCtrlx  KeyPressEvent = KeyPressEvent{gdk.KEY_x, gdk.CONTROL_MASK}
	KeyCtrlf  KeyPressEvent = KeyPressEvent{gdk.KEY_f, gdk.CONTROL_MASK}
	KeyCtrls  KeyPressEvent = KeyPressEvent{gdk.KEY_s, gdk.CONTROL_MASK}
)

// Action Handler applies the Action regardless of what key was pressed and 
// then immediately gives control to the KeyHandler given by the Action.
type ActionHandler struct {
	Action func(*Editor.Editor) KeyHandler
}

func (k *ActionHandler) Handle(e KeyPressEvent, editor *Editor.Editor) (bool, KeyHandler) {
	return true, k.Action(editor)
}

// Returns the default root node
func MakeRoot() KeyHandler {
	root = &rootHandler{&KeyChoice{[]KeyHandler{CtrlXHandler}}}
	return root
}

type rootHandler struct {
	TopLevelChoices KeyHandler
}

func (this *rootHandler) Handle(e KeyPressEvent, editor *Editor.Editor) (bool, KeyHandler) {
	ok, handler := this.TopLevelChoices.Handle(e, editor)
	fmt.Println("rootHandlerProc")
	if ok {
		return ok, handler
	}
	return false, this
}

// GuardHandler checks that a particular key was pressed
// before calling the next keyhandler.
type GuardHandler struct {
	CheckFor KeyPressEvent
	Next     KeyHandler
}

func (this *GuardHandler) Handle(e KeyPressEvent, editor *Editor.Editor) (bool, KeyHandler) {
	if e.Equals(this.CheckFor) {
		fmt.Println("Guard Succeeded")
		_, handler := this.Next.Handle(e, editor)
		return true, handler
	}
	fmt.Println(string(e.KeyVal) + " " + string(e.Modifier))
	fmt.Println("Guard Failed")
	return false, nil
}

// PauseHandler Accepts all input and returns
// its KeyHandler. This means tha the current KeyPressEvent
// is accepted and the next keypressevent is given to Next.
type PauseHandler struct {
	Next KeyHandler
}

func (this *PauseHandler) Handle(e KeyPressEvent, editor *Editor.Editor) (bool, KeyHandler) {
	return true, this.Next
}

// InputHandler waits for input string ended by the Return Key
// and gives it to the Action.
type InputHandler struct {
	Input  string
	Action func(s string, e *Editor.Editor) KeyHandler
}

func (this *InputHandler) Handle(e KeyPressEvent, editor *Editor.Editor) (bool, KeyHandler) {
	if !e.Equals(KeyReturn) {
		this.Input += string(e.KeyVal)
		return true, this
	}
	return true, this.Action(this.Input, editor)
}
