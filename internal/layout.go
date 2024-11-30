package internal

import (
	"fmt"
)

const (
	a4Long  = 297
	a4Short = 210
	spiral  = 12 * posM // Shift from the top
	magic   = 10        // Actually the week and numb font size difference
	posM    = 10        // Multiplication for SVG canvas view.
)

type LayoutType string

const (
	LayoutLeft    LayoutType = "left"
	LayoutRight   LayoutType = "right"
	LayoutBottom  LayoutType = "bottom"
	LayoutTop     LayoutType = "top"
	LayoutSquareV LayoutType = "square_v"
	LayoutSquare  LayoutType = "square"
)

type Layout struct {
	Weekday   Pos
	Number    Pos
	MonthName Pos
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
	indentTop int // for spring perforation
	grid      Pos
	ex        int
}

func (c Canvas) layout(rect Rect, shift Pos, weekGroup int, vertical bool) Layout {
	// Top-left text position
	pos := Pos{
		rect.Pos.X + c.grid.X/2 + shift.X,
		rect.Pos.Y + c.grid.Y/2 + shift.Y,
	}
	numbs := Pos{pos.X, pos.Y + c.grid.Y}

	if vertical {
		pos.Y += c.grid.Y / 8
		numbs = Pos{pos.X + c.grid.X, pos.Y + magic}
	}

	return Layout{
		Weekday:   pos,
		Number:    numbs,
		MonthName: Pos{c.Width / 2, c.Height / 7},
		Shadow:    rect,
		WeekGroup: weekGroup,
		Vertical:  vertical,
	}
}

func NewCanvas(size Size, gridDiv Pos, ex int) Canvas {
	return Canvas{
		SVG: Size{
			Width:  size.Width + ex,
			Height: size.Height + ex,
		},
		Rect: Rect{
			Pos: Pos{0, 0},
			Size: Size{
				Width:  (size.Width + ex) * posM,
				Height: (size.Height + ex) * posM,
			},
		},
		grid: Pos{
			X: size.Width * posM / gridDiv.X,
			Y: size.Height * posM / gridDiv.Y,
		},
		ex: ex * posM,
	}
}

func gridDim(layout LayoutType, album bool) Pos {
	switch layout {
	case LayoutLeft, LayoutRight:
		if album {
			return Pos{20, 15}
		}
		return Pos{15, 22}
	case LayoutBottom, LayoutTop:
		if album {
			return Pos{21, 18}
		}
		return Pos{14, 18}
	case LayoutSquareV:
		if album {
			return Pos{18, 18}
		}
		return Pos{9, 20}
	case LayoutSquare:
		if album {
			return Pos{18, 18}
		}
		return Pos{10, 18}
	}
	panic("Undefined layout")
}

func (c Canvas) Layout(t LayoutType, album bool, pos Pos) Layout {
	switch t {
	case LayoutLeft:
		if album {
			return c.LayoutLeftAlbum()
		}
		return c.LayoutLeftPortrait()
	case LayoutRight:
		if album {
			return c.LayoutRightAlbum()
		}
		return c.LayoutRightPortrait()
	case LayoutBottom:
		if album {
			return c.LayoutBottomAlbum()
		}
		return c.layoutBottomPortrait()
	case LayoutTop:
		if album {
			return c.LayoutTopAlbum()
		}
		return c.LayoutTopPortrait()
	case LayoutSquareV:
		return c.LayoutSquare(pos, true)
	case LayoutSquare:
		return c.LayoutSquare(pos, false)
	}
	panic("Undefined layout")
}

func (c Canvas) LayoutLeftAlbum() Layout {
	ex := c.ex
	width := c.grid.X<<2 + c.grid.X>>3
	l := c.layout(
		Rect{Pos{c.X + ex/2, c.Y + ex/2}, Size{width, c.Height - ex}},
		Pos{Y: spiral},
		2, true)
	l.Shadow.Pos = Pos{c.X, c.Y}
	l.Shadow.Size.Width += ex / 2
	l.Shadow.Size.Height += ex
	return l
}

func (c Canvas) LayoutLeftPortrait() Layout {
	ex := c.ex
	width := c.grid.X*3 + c.grid.X>>3
	l := c.layout(
		Rect{Pos{c.X + ex/2, c.Y + ex/2}, Size{width, c.Height - ex}},
		Pos{Y: spiral},
		3, true)
	l.Shadow.Pos = Pos{c.X, c.Y}
	l.Shadow.Size.Width += ex / 2
	l.Shadow.Size.Height += ex
	return l
}

func (c Canvas) LayoutRightAlbum() Layout {
	l := c.LayoutLeftAlbum()
	l.Shadow.X = c.Width - l.Shadow.Width
	l.Weekday.X += l.Shadow.X + c.grid.X*3 - c.ex/2
	l.Number.X += l.Shadow.X - c.grid.X - c.ex/2
	return l
}

func (c Canvas) LayoutRightPortrait() Layout {
	l := c.LayoutLeftPortrait()
	l.Shadow.X = c.Width - l.Shadow.Width
	l.Weekday.X += l.Shadow.X + c.grid.X*2 - c.ex/2
	l.Number.X += l.Shadow.X - c.grid.X - c.ex/2
	return l
}

func (c Canvas) LayoutBottomAlbum() Layout {
	ex := c.ex
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

func (c Canvas) layoutBottomPortrait() Layout {
	height := c.grid.Y*3 + c.grid.Y>>3
	l := c.layout(
		Rect{
			Pos{c.X + c.ex/2, c.Height - height},
			Size{c.Width, height},
		},
		Pos{},
		2, false,
	)
	l.Shadow.Pos.X -= c.ex / 2
	l.Shadow.Pos.Y -= c.ex / 2
	l.Shadow.Width += c.ex
	l.Shadow.Height += c.ex
	return l
}

func (c Canvas) LayoutTopAlbum() Layout {
	ex := c.ex
	height := c.grid.Y * 3
	l := c.layout(
		Rect{
			Pos{c.X + ex/2, c.Y + ex/2},
			Size{c.Width - ex, height + spiral},
		},
		Pos{magic, spiral},
		3, false,
	)
	l.MonthName.Y = c.Height - c.grid.Y
	l.Shadow.Pos = c.Pos
	l.Shadow.Width += ex
	l.Shadow.Height += ex / 2
	return l
}

func (c Canvas) LayoutTopPortrait() Layout {
	l := c.layout(
		Rect{
			Pos{c.X + c.ex/2, c.Y + c.ex/2},
			Size{c.Width - c.ex, c.grid.Y*3 + spiral}},
		Pos{Y: spiral},
		2, false)

	l.Shadow.Pos = c.Pos
	l.Shadow.Width += c.ex
	l.Shadow.Height += spiral + c.ex/2
	return l
}

func (c Canvas) LayoutSquare(pos Pos, vertical bool) Layout {
	size := Size{c.grid.X * 7, c.grid.Y * 7}
	if pos.X <= 0 {
		pos.X = c.Width/2 - size.Width/2
	}
	if pos.Y <= 0 {
		pos.Y = c.Height/2 - size.Height/2
	}
	return c.layout(Rect{pos, size}, Pos{}, 1, vertical)
}
