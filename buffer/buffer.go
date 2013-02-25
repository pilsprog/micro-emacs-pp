package buffer

import (
	"bufio"
	"bytes"
	"github.com/mattn/go-gtk/gtk"
	"os"
)

func OpenFileInBuffer(tb *gtk.TextBuffer, f string) (err error) {
	var (
		part   []byte
		prefix bool
		start  gtk.TextIter
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

func SaveCurrentOpenFile(tb *gtk.TextBuffer, f string) (err error) {
	var (
		start gtk.TextIter
		end gtk.TextIter
	)

	tb.GetStartIter(&start)
	tb.GetEndIter(&end)

	str := tb.GetText(&start, &end, false)

	fo, err := os.Create(f + "~")
	if err != nil {
		return
	}
	fo.WriteString(str)
	fo.Close()
	return nil
}
