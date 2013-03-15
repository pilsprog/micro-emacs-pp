package Editor
//import "bufio"
import "io"
import "os"
import "github.com/mattn/go-gtk/gtk"	

type Editor struct {
  Filename string
  Buf io.ReadWriter
}

func (e * Editor) OpenFile(f string){
	fo, err := os.Open(f)
	if err != nil {
		return
	}
  io.Copy(fo,e.Buf)
  
}

func (e * Editor) SaveFile(f string){
	fo, err := os.Create(f+"~")
	if err != nil {
		return
	}
  io.Copy(e.Buf,fo)
  
}

type GtkTextBufferReadWriter struct{
  CurrIt gtk.TextIter
  Buf *gtk.TextBuffer
}

func (tbw *GtkTextBufferReadWriter) Read(p []byte) (n int, err error){
  var enditer gtk.TextIter
  tbw.Buf.GetEndIter(&enditer)

  if tbw.CurrIt == enditer {
    return 0,io.EOF
  } else {
    a := []byte(tbw.Buf.GetText(&tbw.CurrIt,&enditer,false))

    tbw.Buf.GetIterAtOffset(&tbw.CurrIt,tbw.CurrIt.GetOffset() + len(a))

    for i := 0; i < len(a);i++ {
      p[i] = a[i] 
    }

    return len(a),nil
  }
  return 0,nil
}
func (tbw *GtkTextBufferReadWriter) Write(p []byte) (n int, err error){

  return 0,nil

}
