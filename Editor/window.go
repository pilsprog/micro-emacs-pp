package editor

import (
	"errors"
)

type Mark int
type Point int

type Direction int

const (
	None Direction = iota
	Horizontal
	Vertical
)

type Window interface {
	GetBuffer() (Buffer, error)
	SetBuffer(Buffer) error
	GetModeline() (Buffer, error)
	GetMark() Mark
	GetPoint() Point
	GetLeft() Window
	GetRight() Window
	GetSplit() Direction
	Split(Direction) error
}
