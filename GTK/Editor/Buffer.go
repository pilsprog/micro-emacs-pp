package Editor

import (
	"fmt"
	"github.com/mattn/go-gtk/gtk"
	"io"
)

type GtkTextBufferReadWriter struct {
	View   *gtk.Widget
	CurrIt gtk.TextIter
	Buf    *gtk.TextBuffer
}

func (this *GtkTextBufferReadWriter) Clear() {
	var start gtk.TextIter
	var end gtk.TextIter
	this.Buf.GetStartIter(&start)
	this.Buf.GetEndIter(&end)
	this.Buf.Delete(&start, &end)
	this.Buf.GetStartIter(&this.CurrIt)
}

func (this *GtkTextBufferReadWriter) SetItStart() {
	this.Buf.GetStartIter(&this.CurrIt)
}

func (this *GtkTextBufferReadWriter) SetItEnd() {
	this.Buf.GetEndIter(&this.CurrIt)
}

func (this *GtkTextBufferReadWriter) GrabFocus() {
	this.View.GrabFocus()
}

func (tbw *GtkTextBufferReadWriter) Read(p []byte) (n int, err error) {
	var enditer gtk.TextIter
	tbw.Buf.GetEndIter(&enditer)

	if tbw.CurrIt == enditer {
		return 0, io.EOF
	}
	a := []byte(tbw.Buf.GetText(&tbw.CurrIt, &enditer, false))

	tbw.Buf.GetIterAtOffset(&tbw.CurrIt, tbw.CurrIt.GetOffset()+len(a))

	for i := 0; i < len(a); i++ {
		p[i] = a[i]
	}

	return len(a), nil
}

func (tbw *GtkTextBufferReadWriter) Write(p []byte) (n int, err error) {
	fmt.Println(string(p))
	tbw.Buf.Insert(&tbw.CurrIt, string(p))
	return len(p), nil
}
