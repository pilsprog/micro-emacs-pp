package Editor

import (
	"io"
)

type Buffer interface {
	io.ReadWriter
	GrabFocus()
	Clear()
	SetItStart()
	SetItEnd()
}
