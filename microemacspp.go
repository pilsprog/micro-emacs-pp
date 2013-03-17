package main

import (
	"os"
	"unsafe"
	"github.com/mattn/go-gtk/gtk"
	"github.com/mattn/go-gtk/gdk"
	"github.com/mattn/go-gtk/glib"
	gsv "github.com/mattn/go-gtk/gtksourceview"
  "micro-emacs-pp/Editor"
  "micro-emacs-pp/KeyHandler"
  "fmt"
)

var textbuffer * gtk.TextBuffer
var sourceview * gsv.SourceView
var textview   * gtk.TextView
var fileName string
var microemacs Editor.Editor
var keyh KeyHandler.KeyHandler = KeyHandler.MakeRoot()

func main() {

	gtk.Init(&os.Args)
	window := gtk.NewWindow(gtk.WINDOW_TOPLEVEL)
	window.SetTitle("µemacs/pp")
	window.Connect("destroy", gtk.MainQuit)
  window.Connect("key-press-event",handleKeyPress);

	swin := gtk.NewScrolledWindow(nil, nil)
	sourcebuffer := gsv.NewSourceBufferWithLanguage(gsv.SourceLanguageManagerGetDefault().GetLanguage("cpp"))
	sourceview = gsv.NewSourceViewWithBuffer(sourcebuffer)

	var start gtk.TextIter
	sourcebuffer.GetStartIter(&start)
	sourcebuffer.Insert(&start, `writing stuff, awww yea!`)

	textview = gtk.NewTextView()
  textbuffer = textview.GetBuffer()

  var bufiter gtk.TextIter
  sourcebuffer.GetStartIter(&bufiter)
  bufWrapper := Editor.GtkTextBufferReadWriter{&sourceview.TextView.Container.Widget,bufiter,&sourcebuffer.TextBuffer}
  var comiter gtk.TextIter
  textbuffer.GetStartIter(&comiter)
  comWrapper := Editor.GtkTextBufferReadWriter{&textview.Container.Widget,comiter,textbuffer}
  microemacs = Editor.Editor{"",&bufWrapper,&comWrapper}

	vbox := gtk.NewVBox(false,0)
	vbox.PackStart(swin, true, true,0)
	vbox.PackEnd(textview,false, true,0)

	swin.Add(sourceview)

	window.Add(vbox)
	window.SetSizeRequest(300, 300)
	window.ShowAll()

	gtk.Main()
}

func handleKeyPress(ctx *glib.CallbackContext){
  arg := ctx.Args(0)
  kev := *(**gdk.EventKey)(unsafe.Pointer(&arg))

  fmt.Println("---------------------\n\n")

  kpe := KeyHandler.KeyPressEvent{int(kev.Keyval),0}
  if (gdk.ModifierType(kev.State) & gdk.CONTROL_MASK) != 0{
    fmt.Println("Added Control Mask")
    kpe.Modifier = gdk.CONTROL_MASK
  }

  _, keyh = keyh.Handle(kpe,&microemacs)
}
