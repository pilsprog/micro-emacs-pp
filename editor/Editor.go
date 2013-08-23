package editor

import (
	"fmt"
	"io"
	"os"
)

type Editor struct {
	Filename   string
	Buf        Buffer
	CommandBuf Buffer
}

func (e *Editor) OpenFile(f string) {
	fo, err := os.Open(f)
	e.Filename = f
	if err != nil {
		fmt.Println("Error Opening File!")
		return
	}
	e.Buf.Clear()
	io.Copy(e.Buf, fo)
}

func (e *Editor) SaveFile(f string) {
	if len(f) == 0 {
		f = e.Filename
	}
	fo, err := os.Create(f + "~")
	if err != nil {
		fmt.Println("Error Saving File!")
		return
	}
	e.Buf.SetItStart()
	io.Copy(fo, e.Buf)
}
