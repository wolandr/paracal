package internal

import (
	"io"
	"strings"
	"time"

	"github.com/ajstarks/svgo"
)

func Draw(w io.Writer, year int, month time.Month, style Style) {
	paper := svg.New(w)
	defer paper.End()

	isAlbum := style.Size.Width > style.Size.Height

	if style.Grid.X == 0 || style.Grid.Y == 0 {
		style.Grid = gridDim(style.Layout, isAlbum)
	}

	canvas := NewCanvas(
		Size{Width: style.Size.Width, Height: style.Size.Height},
		style.Grid,
		style.Extend)

	layout := canvas.Layout(style.Layout, isAlbum)

	if style.MonthNamePos.X > 0 {
		layout.MonthName.X = (style.MonthNamePos.X + style.Extend/2) * posM
	}
	if style.MonthNamePos.Y > 0 {
		layout.MonthName.Y = (style.MonthNamePos.Y + style.Extend/2) * posM
	}

	paper.StartviewUnit(canvas.SVG.Width, canvas.SVG.Height, "mm", canvas.X, canvas.Y, canvas.Width, canvas.Height)

	// Background
	if strings.HasPrefix(style.Background, "#") {
		paper.Rect(0, 0, canvas.Width, canvas.Height, "fill:"+style.Background)
	} else if len(style.Background) > 0 {
		paper.Image(0, 0, canvas.Width, canvas.Height, style.Background, `preserveAspectRatio="none"`)
	}

	// Calendar text background
	paper.Rect(layout.Shadow.X, layout.Shadow.Y, layout.Shadow.Width, layout.Shadow.Height,
		style.ColorStyle.Shadow.style().String())

	// Day of week text group
	grid := Pos{canvas.grid.X, 0}
	if layout.Vertical {
		grid = Pos{0, canvas.grid.Y}
	}
	paper.Group(style.ColorStyle.Weekday.style(), layout.Weekday.transform())
	for i := 0; i < layout.WeekGroup*7; i++ {
		fill := ""
		if i%7 > 4 {
			fill = style.ColorStyle.Weekend.style().String()
		}
		paper.Text(grid.X*i, grid.Y*i, style.WeekdayNames[i%7], fill)
	}
	paper.Gend()

	// Numbers text group
	paper.Group(style.ColorStyle.Number.style(), layout.Number.transform())
	t := backMonday(year, month)
	for i := 0; i < 42; i++ {
		day, fill := style.day(t, month)
		idx := i / (layout.WeekGroup * 7)
		idy := i % (layout.WeekGroup * 7)
		if !layout.Vertical {
			tmp := idx
			idx = idy
			idy = tmp
		}
		paper.Text(canvas.grid.X*idx, canvas.grid.Y*idy, day, fill.String())
		t = t.Add(24 * time.Hour)
	}
	paper.Gend()

	// Month name
	paper.Text(layout.MonthName.X, layout.MonthName.Y, style.MonthName, style.ColorStyle.Month.style())
}

func backMonday(year int, month time.Month) time.Time {
	return backMondayFrom(time.Date(year, month, 1, 0, 0, 0, 0, time.UTC))
}

func backMondayFrom(t time.Time) time.Time {
	for {
		if t.Weekday() == time.Monday {
			return t
		}
		t = t.Add(-24 * time.Hour)
	}
}
