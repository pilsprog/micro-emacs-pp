package KeyHandler

import (
	"github.com/mattn/go-gtk/gdk"
  "micro-emacs-pp/Editor"
  "C"
  "fmt"
)

type KeyPressEvent struct{
  KeyVal int
  Modifier gdk.ModifierType
}

func (k1 KeyPressEvent) Equals(k2 KeyPressEvent) bool{
  return k1.KeyVal == k2.KeyVal && k1.Modifier == k2.Modifier
}

func (k1 KeyPressEvent) toChar() string{
  return "a"
}

type KeyHandler interface{
  Handle(e KeyPressEvent, editor * Editor.Editor) (bool, KeyHandler)
}

type KeyChoice struct{
  choices []KeyHandler
}

func (k *KeyChoice) Handle(e KeyPressEvent,editor * Editor.Editor) (bool,KeyHandler){
  for _,choice := range k.choices {
   ok, handler := choice.Handle(e,editor)
   if ok {
     return ok, handler
   }
  }
  return false,root
}

var(
  root KeyHandler
  CtrlXHandler KeyHandler =
    &GuardHandler{
      KeyCtrlx,
      &PauseHandler{&KeyChoice{[]KeyHandler{
        CtrlFHandler,
        CtrlSHandler}}}}
  CtrlFHandler KeyHandler =
    &GuardHandler{
      KeyCtrlf,
      &ActionHandler{ func(e *Editor.Editor) KeyHandler{
        e.CommandBuf.GrabFocus()
        e.CommandBuf.Clear()
        e.CommandBuf.Write([]byte("Find-file:"))
        return &InputHandler{"", func(s string,e * Editor.Editor) KeyHandler{
          fmt.Println(s)
          e.OpenFile(s)
          e.CommandBuf.Clear()
          e.CommandBuf.Write([]byte("File Opened!"))
          e.Buf.GrabFocus()
          return root }}}}}

  CtrlSHandler KeyHandler =
    &GuardHandler{
      KeyCtrls,
      &ActionHandler{ func(e *Editor.Editor) KeyHandler{
        e.CommandBuf.GrabFocus()
        e.CommandBuf.Clear()
        return &InputHandler{"", func(s string,e * Editor.Editor) KeyHandler{
          e.SaveFile(s)
          e.CommandBuf.Clear()
          e.CommandBuf.Write([]byte("File Saved!"))
          e.Buf.GrabFocus()
          return root }}}}}

  KeyReturn KeyPressEvent = KeyPressEvent{gdk.KEY_Return,0}
  KeyCtrlx KeyPressEvent = KeyPressEvent{gdk.KEY_x,gdk.CONTROL_MASK}
  KeyCtrlf KeyPressEvent = KeyPressEvent{gdk.KEY_f,gdk.CONTROL_MASK}
  KeyCtrls KeyPressEvent = KeyPressEvent{gdk.KEY_s,gdk.CONTROL_MASK}
)

type ActionHandler struct{
   Action func(*Editor.Editor) KeyHandler
}

func (k *ActionHandler) Handle(e KeyPressEvent,editor * Editor.Editor) (bool,KeyHandler){
  return true, k.Action(editor)
}


func MakeRoot() KeyHandler{
   root = &rootHandler{&KeyChoice{[]KeyHandler{CtrlXHandler}}}
   return root
}

type rootHandler struct{
   TopLevelChoices KeyHandler
}

func (this *rootHandler) Handle(e KeyPressEvent,editor * Editor.Editor) (bool,KeyHandler){
  ok,handler := this.TopLevelChoices.Handle(e,editor)
  fmt.Println("rootHandlerProc")
  if ok {
    return ok,handler
  }
  return false,this
}

type GuardHandler struct{
   CheckFor KeyPressEvent
   Next KeyHandler
}


func (this * GuardHandler) Handle(e KeyPressEvent,editor * Editor.Editor) (bool,KeyHandler){
    if e.Equals(this.CheckFor) {
      fmt.Println("Guard Succeeded")
      _,handler := this.Next.Handle(e,editor)
      return true,handler
    }
    fmt.Println(string(e.KeyVal) + " " + string(e.Modifier))
    fmt.Println("Guard Failed")
    return false,nil
}

type PauseHandler struct{
   Next KeyHandler
}

func (this * PauseHandler) Handle(e KeyPressEvent,editor * Editor.Editor) (bool,KeyHandler){
  return true,this.Next
}

type InputHandler struct{
   Input string
   Action func(s string,e *Editor.Editor) KeyHandler
}

func (this * InputHandler) Handle(e KeyPressEvent,editor * Editor.Editor) (bool,KeyHandler){
   if !e.Equals(KeyReturn) {
     this.Input += string(e.KeyVal)
     return true,this
   }
   return true, this.Action(this.Input,editor)
}
