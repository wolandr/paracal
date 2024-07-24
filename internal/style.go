package internal

import (
	"fmt"
	"os"
	"slices"
	"strconv"
	"time"

	"gopkg.in/yaml.v2"
)

type ColorStyle struct {
	Weekday  Font
	Number   Font
	Month    Font
	Shadow   Fill
	Holiday  Fill
	Weekend  Fill
	Shortday Fill
	Ghost    string
}

func (s *ColorStyle) update(a ColorStyle) {
	s.Weekday.update(a.Weekday)
	s.Number.update(a.Number)
	s.Month.update(a.Month)
	s.Shadow.update(a.Shadow)
	s.Holiday.update(a.Holiday)
	s.Weekend.update(a.Weekend)
	s.Shortday.update(a.Shortday)
	replace(&s.Ghost, a.Ghost)
}

type Style struct {
	Size   Size
	Extend int // Extend image width and height in mm after draw.
	Grid   Pos

	Layout       LayoutType
	MonthName    string
	MonthNamePos Pos
	Holidays     []int
	NotWeekend   []int
	ShortDays    []int
	Background   string

	ColorStyle   ColorStyle
	WeekdayNames [7]string
}

func DefaultStyle() Style {
	return Style{
		ColorStyle: ColorStyle{
			Weekday:  Font{Size: "60", Color: "#333333"},
			Number:   Font{Size: "80", Weight: "600", Color: "#333333"},
			Month:    Font{Size: "150", Color: "#449955", Stroke: "gray"},
			Shadow:   Fill{Opacity: "0.5", Color: "white"},
			Holiday:  Fill{Color: "#b03333"},
			Weekend:  Fill{Color: "#b03333"},
			Shortday: Fill{Color: "#584848"},
			Ghost:    "0.4",
		},
		WeekdayNames: weekdays(),

		Size: Size{Width: a4Long, Height: a4Short},
	}
}

func LoadStyle(path string) (*Style, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	style := DefaultStyle()
	if err := yaml.UnmarshalStrict(data, &style); err != nil {
		return nil, err
	}

	return &style, nil
}

func (s Style) day(t time.Time, cal time.Month) (string, SvgStyle) {
	style := make(SvgStyle)
	switch {
	case s.isHoliday(t):
		style = s.ColorStyle.Holiday.style()
	case s.isWeekend(t):
		style = s.ColorStyle.Weekend.style()
	case s.isShortday(t):
		style = s.ColorStyle.Shortday.style()
	}
	if t.Month() != cal {
		style["opacity"] = s.ColorStyle.Ghost
	}
	return strconv.Itoa(t.Day()), style
}

func (s Style) isWeekend(day time.Time) bool {
	if day.Weekday() != time.Saturday && day.Weekday() != time.Sunday {
		return false
	}
	return !slices.Contains(s.NotWeekend, day.Day())
}

func (s Style) isHoliday(day time.Time) bool {
	return slices.Contains(s.Holidays, day.Day())
}

func (s Style) isShortday(day time.Time) bool {
	return slices.Contains(s.ShortDays, day.Day())
}

func (s Style) String() string {
	dump, err := yaml.Marshal(s)
	if err != nil {
		return fmt.Sprintf("Dump config error: %s\n", err)
	}
	return string(dump)
}

type SvgStyle map[string]string

func (s SvgStyle) String() string {
	style := ""
	for k, v := range s {
		if v != "" {
			style += k + ":" + v + ";"
		}
	}
	return style
}

type Fill struct {
	Opacity string
	Color   string
}

func (f *Fill) update(a Fill) {
	replace(&f.Opacity, a.Opacity)
	replace(&f.Color, a.Color)
}

func (f Fill) style() SvgStyle {
	return SvgStyle{"opacity": f.Opacity, "fill": f.Color}
}

type Font struct {
	Family      string
	Size        string
	Weight      string // bold
	Color       string
	Stroke      string
	StrokeWidth string
}

func (f *Font) update(a Font) {
	replace(&f.Family, a.Family)
	replace(&f.Size, a.Size)
	replace(&f.Weight, a.Weight)
	replace(&f.Color, a.Color)
	replace(&f.Stroke, a.Stroke)
	replace(&f.StrokeWidth, a.StrokeWidth)
}

func (f Font) style() string {
	return SvgStyle{
		"text-align": "center", "fill-opacity": "1", "fill-rule": "nonzero", "text-anchor": "middle",
		"font-family":  f.Family,
		"font-size":    f.Size,
		"font-weight":  f.Weight,
		"fill":         f.Color,
		"stroke":       f.Stroke,
		"stroke-width": f.StrokeWidth,
	}.String()
}

func weekdays() [7]string {
	weekdays := [7]string{}
	for i := 0; i < 7; i++ {
		weekdays[i] = time.Weekday(i).String()[:2]
	}
	return weekdays
}

func replace(s *string, with string) {
	if len(with) > 0 {
		*s = with
	}
}
