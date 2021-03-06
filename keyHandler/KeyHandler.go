// Package Keyhandler provides combinators for
// writing corded keyevents with calls to microemacspp/Editor. The combinators
// builds a tree of possible key combinations.
package keyhandler

import (
	"C"
	"fmt"
	"micro-emacs-pp/editor"
)

var (
	Root KeyHandler
	KeyReturn KeyPressEvent
)

// Returns the default root node
func MakeRoot(khs ...KeyHandler) KeyHandler {
	Root = &rootHandler{KeyChoice(khs...)}
	return Root
}

func SetKeyReturn(k KeyPressEvent) KeyPressEvent {
	KeyReturn = k
	return KeyReturn
}

type Modifier int

const (
	NONE Modifier = iota
	CTRL
	FN
	HYPER
	META
	SUPER
)

func (m Modifier) String() string {
	switch {
	case m == CTRL:
		return "Ctrl"
	case m == FN:
		return "Fn"
	case m == HYPER:
		return "Hyper"
	case m == META:
		return "Meta"
	case m == SUPER:
		return "Super"
	}
	return "None"
}

type KeyPressEvent interface {
	GetKeyValue() int
	GetModifier() Modifier
	Equals(KeyPressEvent) bool
}

// KeyHandler is the interface that every node in the tree implements. Accepts
// returns true if the KeyHandler takes resposibility for the given keypresses and
// Insert attempts to insert a new KeyHandler and returns true if it is
// successful.
type KeyHandler interface {
	Handle(e KeyPressEvent, editor *editor.Editor) (bool, KeyHandler)
	Accepts(e []KeyPressEvent) bool
	Insert(e []KeyPressEvent, h KeyHandler) bool
	//  Replace(e []KeyPressEvent,h KeyHandler) bool
}

// KeyChoice is a KeyHandler which represents a choice between several different
// keyHandlers. KeyChoice calls every KeyHandler in its list until one succeeds.
func KeyChoice(choices ...KeyHandler) KeyHandler {
	return &keyChoice{[]KeyHandler(choices)}
}

type keyChoice struct {
	choices []KeyHandler
}

func (k *keyChoice) Handle(e KeyPressEvent, editor *editor.Editor) (bool, KeyHandler) {
	for _, choice := range k.choices {
		ok, handler := choice.Handle(e, editor)
		if ok {
			return ok, handler
		}
	}
	return false, Root
}

func (k *keyChoice) Accepts(e []KeyPressEvent) bool {
	for _, choice := range k.choices {
		if choice.Accepts(e) {
			return true
		}
	}
	return false
}

func (k *keyChoice) Insert(e []KeyPressEvent, kh KeyHandler) bool {
	fmt.Println("keyChoice insert")
	for _, choice := range k.choices {
		accepts := choice.Accepts(e)
		if accepts {
			return false
		}

		inserted := choice.Insert(e, kh)
		if inserted {
			fmt.Println("insert succeeded in child")
			return true
		}
	}
	fmt.Println("Appended new choice")
	k.choices = append(k.choices, makeGuards(e, kh))
	return false
}

func makeGuards(e []KeyPressEvent, kh KeyHandler) KeyHandler {
	if len(e) == 0 {
		return kh
	}
	return GuardHandler(e[0], makeGuards(e[1:], kh))
}

// Action Handler applies the Action regardless of what key was pressed and
// then immediately gives control to the KeyHandler given by the Action.
func ActionHandler(Action func(*editor.Editor) KeyHandler) KeyHandler {
	return &actionHandler{Action}
}

type actionHandler struct {
	action func(*editor.Editor) KeyHandler
}

func (k *actionHandler) Handle(e KeyPressEvent, editor *editor.Editor) (bool, KeyHandler) {
	return true, k.action(editor)
}

func (k *actionHandler) Accepts(e []KeyPressEvent) bool {
	return true
}

func (k *actionHandler) Insert(e []KeyPressEvent, kh KeyHandler) bool {
	return false
}

type rootHandler struct {
	TopLevelChoices KeyHandler
}

func (this *rootHandler) Handle(e KeyPressEvent, editor *editor.Editor) (bool, KeyHandler) {
	ok, handler := this.TopLevelChoices.Handle(e, editor)
	if ok {
		return ok, handler
	}
	return true, this
}

func (this *rootHandler) Accepts(e []KeyPressEvent) bool {
	return true
}

func (this *rootHandler) Insert(e []KeyPressEvent, kh KeyHandler) bool {
	fmt.Println("rootInsert")
	return this.TopLevelChoices.Insert(e, kh)
}

// GuardHandler checks that a particular key was pressed
// before calling the next keyhandler.
func GuardHandler(chk KeyPressEvent, next KeyHandler) KeyHandler {
	return &guardHandler{chk, next}
}

type guardHandler struct {
	checkFor KeyPressEvent
	next     KeyHandler
}

func (this *guardHandler) Handle(e KeyPressEvent, editor *editor.Editor) (bool, KeyHandler) {
	if e.Equals(this.checkFor) {
		fmt.Println("Guard Succeeded")
		_, handler := this.next.Handle(e, editor)
		return true, handler
	}
	fmt.Println(string(e.GetKeyValue()), " : ", e.GetModifier())
	fmt.Println("Guard Failed")
	return false, nil
}

func (this *guardHandler) Accepts(e []KeyPressEvent) bool {
	if e[0].Equals(this.checkFor) {
		return this.next.Accepts(e[1:])
	}
	return false
}

func (this *guardHandler) Insert(e []KeyPressEvent, kh KeyHandler) bool {
	if e[0].Equals(this.checkFor) {
		if this.next.Accepts(e[1:]) {
			return this.next.Insert(e[1:], kh)
		} else {
			this.next = KeyChoice(this.next, makeGuards(e[1:], kh))
		}
		return true
	}
	return false
}

// PauseHandler Accepts all input and returns
// its KeyHandler. This means tha the current KeyPressEvent
// is accepted and the next keypressevent is given to Next.
func PauseHandler(Next KeyHandler) KeyHandler {
	return &pauseHandler{Next}
}

type pauseHandler struct {
	next KeyHandler
}

func (this *pauseHandler) Handle(e KeyPressEvent, editor *editor.Editor) (bool, KeyHandler) {
	return true, this.next
}

func (this *pauseHandler) Accepts(e []KeyPressEvent) bool {
	return this.next.Accepts(e)
}

func (this *pauseHandler) Insert(e []KeyPressEvent, kh KeyHandler) bool {
	return this.next.Insert(e, kh)
}

// InputHandler waits for input string ended by the Return Key
// and gives it to the Action.
func InputHandler(action func(s string, e *editor.Editor) KeyHandler) KeyHandler {
	return &inputHandler{action}
}

type inputHandler struct {
	action func(s string, e *editor.Editor) KeyHandler
}

func (this *inputHandler) Handle(e KeyPressEvent, editor *editor.Editor) (bool, KeyHandler) {
	if !e.Equals(KeyReturn) {
		return true, this
	}
	buff := make([]byte, 512)
	editor.CommandBuf.SetItStart()
	n, _ := editor.CommandBuf.Read(buff)
	return true, this.action(string(buff[0:n]), editor)
}

func (this *inputHandler) Accepts(e []KeyPressEvent) bool {
	return true
}

func (this *inputHandler) Insert(e []KeyPressEvent, kh KeyHandler) bool {
	return false
}
