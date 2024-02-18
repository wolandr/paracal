package internal

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/ajstarks/svgo"
)

func Draw(w io.Writer, year int, month time.Month, layout LayoutType, style Style, background string, ex int) {
	paper := svg.New(w)
	defer paper.End()

	var monthPos Pos
	if s, ok := style.MonthStyle[int(month)]; ok {
		if s.Layout != "" {
			layout = s.Layout
		}
		monthPos = s.MonthPos
	}

	canvas := NewAlbumCanvas(ex)
	pos := canvas.Layout(layout, ex)
	if monthPos.X > 0 {
		pos.Month.X = monthPos.X
	}
	if monthPos.Y > 0 {
		pos.Month.Y = monthPos.Y
	}

	fmt.Printf("Month label position: %d %d\n", pos.Month.X, pos.Month.Y)
	paper.StartviewUnit(canvas.SVG.Width, canvas.SVG.Height, "mm", canvas.X, canvas.Y, canvas.Width, canvas.Height)

	// Background
	if strings.HasPrefix(background, "#") {
		paper.Rect(0, 0, canvas.Width, canvas.Height, "fill:"+background)
	} else if len(background) > 0 {
		paper.Image(0, 0, canvas.Width, canvas.Height, background, `preserveAspectRatio="none"`)
	}

	// Calendar text background
	paper.Rect(pos.Shadow.X, pos.Shadow.Y, pos.Shadow.Width, pos.Shadow.Height, style.Shadow.style().String())

	// Day of week text group
	grid := Pos{canvas.grid.X, 0}
	if pos.Vertical {
		grid = Pos{0, canvas.grid.Y}
	}
	paper.Group(style.Weekday.style(), pos.Weekday.transform())
	for i := 0; i < pos.WeekGroup*7; i++ {
		fill := ""
		if i%7 > 4 {
			fill = style.Weekend.style().String()
		}
		paper.Text(grid.X*i, grid.Y*i, style.Weekdays[i%7], fill)
	}
	paper.Gend()

	// Numbers text group
	paper.Group(style.Number.style(), pos.Number.transform())
	t := backMonday(year, month)
	for i := 0; i < 42; i++ {
		day, fill := style.day(t, month)
		idx := i / (pos.WeekGroup * 7)
		idy := i % (pos.WeekGroup * 7)
		if !pos.Vertical {
			tmp := idx
			idx = idy
			idy = tmp
		}
		paper.Text(canvas.grid.X*idx, canvas.grid.Y*idy, day, fill.String())
		t = t.Add(24 * time.Hour)
	}
	paper.Gend()

	// Month name
	paper.Text(pos.Month.X, pos.Month.Y, style.Months[month-1], style.Month.style())

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
