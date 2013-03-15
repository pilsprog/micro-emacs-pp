package KeyHandler

import (
	"github.com/mattn/go-gtk/gdk"	
  "microemacspp/Editor"
  "C"
)

type KeyPressEvent struct{
  keyVal C.unsigned
  modifier gdk.ModifierType
}

func (k1 KeyPressEvent) Equals(k2 KeyPressEvent) bool{
  return k1.keyVal == k2.keyVal && k1.modifier == k2.modifier
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

func (k KeyChoice) Handle(e KeyPressEvent,editor * Editor.Editor) (bool,KeyHandler){
   return false, nil 
}

var(
  Root KeyChoice = KeyChoice{[]KeyHandler{CtrlXHandler}}
  CtrlXHandler KeyChoice = KeyChoice{[]KeyHandler{&CtrlFHandler{""},&CtrlSHandler{""}}}
  KeyReturn KeyPressEvent = KeyPressEvent{gdk.KEY_Return,0}
)

type CtrlFHandler struct{
  filename string
}

func (h * CtrlFHandler) Handle(e KeyPressEvent,editor * Editor.Editor) (bool,KeyHandler){
   if !e.Equals(KeyReturn) {
     h.filename += e.toChar()
     return true,h
   } else {
     editor.OpenFile(h.filename)
     return true,Root
   }
   return false,nil
}

type CtrlSHandler struct{
  filename string
}

func (h * CtrlSHandler) Handle(e KeyPressEvent,editor * Editor.Editor) (bool,KeyHandler){
  editor.SaveFile(h.filename)
  return true,Root
}

