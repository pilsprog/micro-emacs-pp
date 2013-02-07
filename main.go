package main

import (
	"os"
	"unsafe"	
	"github.com/mattn/go-gtk/gtk"	
	"github.com/mattn/go-gtk/gdk"	
	"github.com/mattn/go-gtk/glib"	
	gsv "github.com/mattn/go-gtk/gtksourceview"	
)

var textbuffer * gtk.TextBuffer
var sourceview * gsv.SourceView
var textview   * gtk.TextView
var fileName string

var readingFileName bool

func main() {
  readingFileName = false

	gtk.Init(&os.Args)
	window := gtk.NewWindow(gtk.WINDOW_TOPLEVEL)
	window.SetTitle("microemacspp")
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
	var start gtk.TextIter
	var end gtk.TextIter
	textbuffer.GetStartIter(&start)
	textbuffer.GetEndIter(&end)
  if ((gdk.ModifierType(kev.State) & gdk.CONTROL_MASK) != 0) && kev.Keyval == gdk.KEY_x {
    sourceview.SetEditable(false)
	  textbuffer.Insert(&end, "^X")
  } else if kev.Keyval == gdk.KEY_f && !sourceview.GetEditable() && !readingFileName{
    textview.GrabFocus()
	  textbuffer.GetStartIter(&start)
	  textbuffer.GetEndIter(&end)
    textbuffer.Delete(&start,&end)
    textbuffer.Insert(&start,"Find-file: ")
    readingFileName = true
  } else if kev.Keyval == gdk.KEY_s && !sourceview.GetEditable() && !readingFileName && (gdk.ModifierType(kev.State) & gdk.CONTROL_MASK) != 0{
	  textbuffer.GetStartIter(&start)
	  textbuffer.GetEndIter(&end)
    textbuffer.Delete(&start,&end)
    textbuffer.Insert(&start,"Find-file: ")
    sourceview.SetEditable(true)
    textview.SetEditable(false)
  } else if readingFileName {
    if kev.Keyval != gdk.KEY_Return {
      textview.SetEditable(true)
    } else { 
	    var fileStart gtk.TextIter
      textbuffer.GetIterAtOffset(&fileStart,11)
	    textbuffer.GetStartIter(&start)
	    textbuffer.GetEndIter(&end)

      fileName = textbuffer.GetText(&fileStart,&end,true)
      textview.SetEditable(false)
      textbuffer.Delete(&start,&end)
      readingFileName = false
    }
  } else {
	  textbuffer.GetStartIter(&start)
	  textbuffer.GetEndIter(&end)
    textbuffer.Delete(&start,&end)
    textbuffer.Insert(&start,"Find-file: ")
    sourceview.SetEditable(true)
    textview.SetEditable(false)
  }
}
