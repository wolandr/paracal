package internal

import (
	"fmt"
)

const (
	a4Long  = 297
	a4Short = 210
	spiral  = 120
	magic   = 10 // Actually the week and numb font size difference
)

type LayoutType string

const (
	LayoutLeft LayoutType = "left"
	LayoutRight LayoutType = "right"
	LayoutBottom LayoutType = "bottom"
	LayoutTop LayoutType = "top"
	LayoutSquare LayoutType = "square"
	LayoutSquareHor LayoutType = "square_h"
)

type Layout struct {
	Weekday   Pos
	Number    Pos
	Month     Pos
	Shadow    Rect
	WeekGroup int
	Vertical  bool
}

type Pos struct {
	X int
	Y int
}

func (p Pos) transform() string {
	return fmt.Sprintf(`transform="translate(%d,%d)"`, p.X, p.Y)
}

type Size struct {
	Width  int
	Height int
}

type Rect struct {
	Pos
	Size
}

type Canvas struct {
	SVG Size
	Rect
	grid Pos
}

func (c Canvas) layout(rect Rect, shift Pos, weekGroup int, vertical bool) Layout {
	pos := Pos{rect.Pos.X + c.grid.X/2 + shift.X, rect.Pos.Y + c.grid.Y/2 + shift.Y}
	var numbs Pos
	if vertical {
		pos.Y += c.grid.Y / 8
		numbs = Pos{pos.X + c.grid.X, pos.Y + magic}
	} else {
		numbs = Pos{pos.X, pos.Y + c.grid.Y}
	}

	return Layout{
		Weekday:   pos,
		Number:    numbs,
		Month:     Pos{c.Width / 2, c.Height / 7},
		Shadow:    rect,
		WeekGroup: weekGroup,
		Vertical:  vertical,
	}
}

func (c Canvas) LayoutSquare(vertical bool) Layout {
	size := Size{c.grid.X * 7, c.grid.Y * 7}
	pos := Pos{c.Width/2 - size.Width/2, c.Height/2 - size.Height/2}
	return c.layout(Rect{pos, size}, Pos{}, 1, vertical)
}

type AlbumCanvas struct{ Canvas }

func NewAlbumCanvas() AlbumCanvas {
	return AlbumCanvas{Canvas{
		SVG:  Size{a4Long, a4Short},
		Rect: Rect{Pos: Pos{0, 0}, Size: Size{a4Long * 10, a4Short * 10}},
		grid: Pos{a4Short * 10 / 15, a4Short * 10 / 15},
	}}
}

func (c AlbumCanvas) Layout(t LayoutType) Layout {
	switch t {
	case LayoutLeft:
		return c.LayoutLeft()
	case LayoutRight:
		return c.LayoutRight()
	case LayoutBottom:
		return c.LayoutBottom()
	case LayoutTop:
		return c.LayoutTop()
	case LayoutSquare:
		return c.LayoutSquare(true)
	case LayoutSquareHor:
		return c.LayoutSquare(false)
	}
	panic("Undefined layout")
}

func (c AlbumCanvas) LayoutLeft() Layout {
	width := c.grid.X<<2 + c.grid.X>>3
	l := c.layout(Rect{Pos{}, Size{width, c.Height}}, Pos{Y: spiral}, 2, true)
	l.Month.X += width / 2
	return l
}

func (c AlbumCanvas) LayoutRight() Layout {
	l := c.LayoutLeft()
	l.Shadow.X = c.Width - l.Shadow.Width
	l.Weekday.X += l.Shadow.X + c.grid.X*3
	l.Number.X += l.Shadow.X - c.grid.X
	l.Month.X -= l.Shadow.Width
	return l
}

func (c AlbumCanvas) LayoutBottom() Layout {
	height := c.grid.Y*3 + magic
	return c.layout(Rect{Pos{0, c.Height - height}, Size{c.Width, height}}, Pos{magic, magic * 2}, 3, false)
}

func (c AlbumCanvas) LayoutTop() Layout {
	height := c.grid.Y * 3
	l := c.layout(Rect{Pos{0, 0}, Size{c.Width, height + spiral}}, Pos{magic, spiral}, 3, false)
	l.Month.Y = c.Height - c.grid.Y
	return l
}

type PortraitCanvas struct{ Canvas }

func NewPortraitCanvas() PortraitCanvas {
	return PortraitCanvas{Canvas{
		SVG:  Size{a4Short, a4Long},
		Rect: Rect{Pos: Pos{0, 0}, Size: Size{a4Short * 10, a4Long * 10}},
		grid: Pos{a4Short * 10 / 15, a4Short * 10 / 15},
	}}
}

func (c PortraitCanvas) layoutBottom() Layout {
	fat := c.grid.Y*4 - magic
	l := c.layout(Rect{Pos{Y: c.Height - fat}, Size{c.Width, fat}}, Pos{X: 45}, 2, false)
	return l
}
