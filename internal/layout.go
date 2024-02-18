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
	LayoutLeft      LayoutType = "left"
	LayoutRight     LayoutType = "right"
	LayoutBottom    LayoutType = "bottom"
	LayoutTop       LayoutType = "top"
	LayoutSquare    LayoutType = "square"
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

func NewAlbumCanvas(ex int) AlbumCanvas {
	return AlbumCanvas{Canvas{
		SVG:  Size{a4Long + ex, a4Short + ex},
		Rect: Rect{Pos: Pos{0, 0}, Size: Size{(a4Long + ex) * 10, (a4Short + ex) * 10}},
		grid: Pos{a4Short * 10 / 15, a4Short * 10 / 15},
	}}
}

func (c AlbumCanvas) Layout(t LayoutType, ex int) Layout {
	ex = ex * 10
	switch t {
	case LayoutLeft:
		return c.LayoutLeft(ex)
	case LayoutRight:
		return c.LayoutRight(ex)
	case LayoutBottom:
		return c.LayoutBottom(ex)
	case LayoutTop:
		return c.LayoutTop(ex)
	case LayoutSquare:
		return c.LayoutSquare(true)
	case LayoutSquareHor:
		return c.LayoutSquare(false)
	}
	panic("Undefined layout")
}

func (c AlbumCanvas) LayoutLeft(ex int) Layout {
	width := c.grid.X<<2 + c.grid.X>>3
	l := c.layout(Rect{Pos{c.X + ex/2, c.Y + ex/2}, Size{width, c.Height - ex}}, Pos{Y: spiral}, 2, true)
	l.Month.X += width / 2
	l.Month.Y += ex / 2
	l.Shadow.Pos = Pos{c.X, c.Y}
	l.Shadow.Size.Width += ex / 2
	l.Shadow.Size.Height += ex
	return l
}

func (c AlbumCanvas) LayoutRight(ex int) Layout {
	l := c.LayoutLeft(ex)
	l.Shadow.X = c.Width - l.Shadow.Width
	l.Weekday.X += l.Shadow.X + c.grid.X*3 - ex/2
	l.Number.X += l.Shadow.X - c.grid.X - ex/2
	l.Month.X -= l.Shadow.Width
	return l
}

func (c AlbumCanvas) LayoutBottom(ex int) Layout {
	height := c.grid.Y*3 + magic
	l := c.layout(
		Rect{
			Pos{c.X + ex/2, c.Y + c.Height - ex/2 - height},
			Size{c.Width - ex, height},
		},
		Pos{magic, magic * 2},
		3, false)
	l.Shadow.X = c.X
	l.Shadow.Width += ex
	l.Shadow.Height += ex / 2
	return l
}

func (c AlbumCanvas) LayoutTop(ex int) Layout {
	height := c.grid.Y * 3
	l := c.layout(Rect{Pos{c.X + ex/2, c.Y + ex/2}, Size{c.Width - ex, height + spiral}}, Pos{magic, spiral}, 3, false)
	l.Month.Y = c.Height - c.grid.Y
	l.Shadow.Pos = c.Pos
	l.Shadow.Width += ex
	l.Shadow.Height += ex / 2
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
