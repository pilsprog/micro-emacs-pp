package main

import (
	"bufio"
	"bytes"
	"os"
	"github.com/mattn/go-gtk/gtk"	
	gsv "github.com/mattn/go-gtk/gtksourceview"	
)

func OpenFileInBuffer(tb *gsv.SourceBuffer,f string) (err error) {
	var (
		part []byte
		prefix bool
		start gtk.TextIter
	)

	file, err := os.Open(f)
	if err != nil {
		return
	}	
	
	tb.GetStartIter(&start)

	reader := bufio.NewReader(file)
	buffer := bytes.NewBuffer(make([]byte, 1024))
	
	for {
		if part, prefix, err = reader.ReadLine(); err != nil {
			return err
		}
		buffer.Write(part)
		if !prefix {
			tb.Insert(&start, buffer.String()+"\n")
			buffer.Reset()
		}
	}
	file.Close()
	return nil
}
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
