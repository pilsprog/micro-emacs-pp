package main

import (
	"github.com/mattn/go-gtk/gdk"
	"github.com/mattn/go-gtk/glib"
	"github.com/mattn/go-gtk/gtk"
	gsv "github.com/mattn/go-gtk/gtksourceview"
	"micro-emacs-pp/editor"
	ed "github.com/pilsprog/micro-emacs-pp/gtk/editor"
	. "github.com/pilsprog/micro-emacs-pp/keyhandler"
	kh "github.com/pilsprog/micro-emacs-pp/gtk/keyhandler"
	"os"
	"unsafe"
)

var (
	textbuffer *gtk.TextBuffer
	sourceview *gsv.SourceView
	textview   *gtk.TextView
	fileName   string
	microemacs editor.Editor
	keyh       KeyHandler = MakeRoot(kh.CtrlXHandler)
)

func main() {

	gtk.Init(&os.Args)
	window := gtk.NewWindow(gtk.WINDOW_TOPLEVEL)
	window.SetTitle("µemacs/pp")
	window.Connect("destroy", gtk.MainQuit)
	window.Connect("key-press-event", handleKeyPress)

	SetKeyReturn(kh.GTKKeyPressEvent{gdk.KEY_Return, 0})

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
	bufWrapper := ed.GtkTextBufferReadWriter{&sourceview.TextView.Container.Widget, bufiter, &sourcebuffer.TextBuffer}
	var comiter gtk.TextIter
	textbuffer.GetStartIter(&comiter)
	comWrapper := ed.GtkTextBufferReadWriter{&textview.Container.Widget, comiter, textbuffer}
	microemacs = editor.Editor{"", &bufWrapper, &comWrapper}

	vbox := gtk.NewVBox(false, 0)
	vbox.PackStart(swin, true, true, 0)
	vbox.PackEnd(textview, false, true, 0)

	swin.Add(sourceview)

	window.Add(vbox)
	window.SetSizeRequest(300, 300)
	window.ShowAll()

	gtk.Main()
}

func handleKeyPress(ctx *glib.CallbackContext) {
	arg := ctx.Args(0)
	kev := *(**gdk.EventKey)(unsafe.Pointer(&arg))

	kpe := kh.GTKKeyPressEvent{int(kev.Keyval), 0}
	if (gdk.ModifierType(kev.State) & gdk.CONTROL_MASK) != 0 {
		kpe.Modifier = gdk.CONTROL_MASK
	}

	_, keyh = keyh.Handle(kpe, &microemacs)
}
