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

	style.Pos.X *= posM
	style.Pos.Y *= posM
	layout := canvas.Layout(style.Layout, isAlbum, style.Pos)

	//if style.MonthNamePos.X > 0 {
	//	layout.MonthName.X = (style.MonthNamePos.X + style.Extend/2) * posM
	//}
	//if style.MonthNamePos.Y > 0 {
	//	layout.MonthName.Y = (style.MonthNamePos.Y + style.Extend/2) * posM
	//}

	paper.StartviewUnit(canvas.SVG.Width, canvas.SVG.Height, "mm", canvas.X, canvas.Y, canvas.Width, canvas.Height)

	// Background
	if strings.HasPrefix(style.Background, "#") {
		paper.Rect(0, 0, canvas.Width, canvas.Height, "fill:"+style.Background)
	} else if len(style.Background) > 0 {
		paper.Image(0, 0, canvas.Width, canvas.Height, style.Background, `preserveAspectRatio="none"`)
	}

	// Calendar text background
	daysN := daysCount(year, month, layout, style.Layout)
	shadowWidth, shadowHeight := layout.Shadow.Width, layout.Shadow.Height
	switch style.Layout {
	case LayoutSquare:
		shadowHeight = shadowHeight * ((daysN / 7) + 1) / 7
	case LayoutSquareV:
		shadowWidth = shadowWidth * ((daysN / 7) + 1) / 7
	}
	paper.Rect(layout.Shadow.X, layout.Shadow.Y, shadowWidth, shadowHeight,
		style.ColorStyle.Shadow.String())

	// Day of week text group
	grid := Pos{canvas.grid.X, 0}
	if layout.Vertical {
		grid = Pos{0, canvas.grid.Y}
	}
	paper.Group(style.ColorStyle.Weekday.String(), layout.Weekday.transform())
	for i := 0; i < layout.WeekGroup*7; i++ {
		fill := ""
		if i%7 > 4 {
			fill = style.ColorStyle.Weekend.String()
		}
		paper.Text(grid.X*i, grid.Y*i, style.WeekdayNames[i%7], fill)
	}
	paper.Gend()

	// Numbers text group
	paper.Group(style.ColorStyle.Number.String(), layout.Number.transform())
	t := backMonday(year, month)
	for i := 0; i < daysN; i++ {
		day, attr := style.day(t, month)
		idx := i / (layout.WeekGroup * 7)
		idy := i % (layout.WeekGroup * 7)
		if !layout.Vertical {
			tmp := idx
			idx = idy
			idy = tmp
		}
		paper.Text(canvas.grid.X*idx, canvas.grid.Y*idy, day, attr.String())
		t = t.Add(24 * time.Hour)
	}
	paper.Gend()

	// Text
	for _, txt := range style.Text {
		txt.X *= posM
		txt.Y *= posM
		paper.Text(txt.X, txt.Y, txt.Title, txt.Style.String())
	}
}

func daysCount(year int, month time.Month, layout Layout, ltype LayoutType) int {
	if ltype != LayoutSquare && ltype != LayoutSquareV {
		return 7 * 6
	}

	t := backMonday(year, month)
	for i := 0; i < 7*6; i++ {
		idx := i / (layout.WeekGroup * 7)
		idy := i % (layout.WeekGroup * 7)
		if !layout.Vertical {
			tmp := idx
			idx = idy
			idy = tmp
		}
		outOfMonth := year*100+int(month) < t.Year()*100+int(t.Month())
		if outOfMonth {
			if (ltype == LayoutSquare && idx == 0) || (ltype == LayoutSquareV && idy == 0) {
				return 7 * 5
			}
		}
		t = t.Add(24 * time.Hour)
	}

	return 7 * 6
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
