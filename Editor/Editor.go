package Editor
import "bufio"
import "io"
import "github.com/mattn/go-gtk/gtk"	

type Editor struct {
  filename string
  buf ReaderWriter
}

func (e * Editor) OpenFile(f string){
	fo, err := os.Open(f)
	if err != nil {
		return
	}
  io.Copy(fo,e.buf)
  
}

func (e * Editor) SaveFile(f string){
	fo, err := os.Create(f+"~")
	if err != nil {
		return
	}
  io.Copy(e.buf,fo)
  
}

type GtkTextBufferReaderWriter{
  currIt gtk.TextIter
  buf gtk.TextBuffer
}

func (tbw *GtkTextBufferWrapper) Read(p []byte) (n int, err error){
  var enditer gtk.textiter
  tbw.buf.getenditer(enditer)

  if tbw.currIt == enditer {
    return 0,EOF
  } else {
    a := []byte(tbw.buf.gettext(&tbw.currit))

    tbw.buf.getiteratoffset(&tbw.currit,tbw.currit.getoffset + len(a))

    for i := 0; i < len(a);i++ {
      p[i] = a[i] 
    }

    return len(a),nil
  }
}
func (tbw *GtkTextBufferWrapper) Write(p []byte) (n int, err error){

  return 0,nil

}
