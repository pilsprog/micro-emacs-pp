package main

import (
	"os"
	"unsafe"	
	"github.com/mattn/go-gtk/gtk"	
	"github.com/mattn/go-gtk/gdk"	
	"github.com/mattn/go-gtk/glib"	
	gsv "github.com/mattn/go-gtk/gtksourceview"	
  "microemacspp/Editor"
  "microemacspp/KeyHandler"
)

var textbuffer * gtk.TextBuffer
var sourceview * gsv.SourceView
var textview   * gtk.TextView
var fileName string
var microemacs Editor.Editor
var keyh       KeyHandler.KeyHandler

var readingFileName bool

func main() {
  readingFileName = false

  k := KeyHandler.Root
  k.Handle(KeyHandler.KeyReturn,nil)
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
  var iter gtk.TextIter
  textbuffer.GetStartIter(&iter)
  wrapper := Editor.GtkTextBufferReadWriter{iter,textbuffer}
  microemacs = Editor.Editor{"",&wrapper}

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

  kpe := KeyHandler.KeyPressEvent{int(kev.Keyval),0}
  if (gdk.ModifierType(kev.State) & gdk.CONTROL_MASK) != 0{
    kpe.Modifier = gdk.CONTROL_MASK
  }

  _, keyh = keyh.Handle(kpe,&microemacs)
}
