package main

import (
	"os"
	"github.com/mattn/go-gtk/gtk"	
	gsv "github.com/mattn/go-gtk/gtksourceview"	
)

func main() {
	gtk.Init(&os.Args)
	window := gtk.NewWindow(gtk.WINDOW_TOPLEVEL)
	window.SetTitle("µemacs/pp")
	window.Connect("destroy", gtk.MainQuit)

	swin := gtk.NewScrolledWindow(nil, nil)
	sourcebuffer := gsv.NewSourceBufferWithLanguage(gsv.SourceLanguageManagerGetDefault().GetLanguage("cpp"))
	sourceview := gsv.NewSourceViewWithBuffer(sourcebuffer)

	var start gtk.TextIter
	sourcebuffer.GetStartIter(&start)
	sourcebuffer.Insert(&start, `writing stuff, awww yea!`)

	swin.Add(sourceview)
	window.Add(swin)
	window.SetSizeRequest(200, 200)
	window.ShowAll()

	gtk.Main()
}
