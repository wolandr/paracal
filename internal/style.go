package internal

import (
	"os"
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

type MonthStyle struct {
	ColorStyle
	Layout   LayoutType
	MonthPos Pos
}

type Style struct {
	ColorStyle
	Weekdays   [7]string
	Months     [12]string
	Holidays   map[int]map[int][]int
	NotWeekend map[int]map[int][]int
	Shortdays  map[int]map[int][]int
	MonthStyle map[int]MonthStyle
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
		Weekdays: weekdays(),
		Months:   months(),
	}
}

func LoadStyle(path string, month time.Month) (*Style, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	style := DefaultStyle()
	if err := yaml.UnmarshalStrict(data, &style); err != nil {
		return nil, err
	}

	if m, found := style.MonthStyle[int(month)]; found {
		style.ColorStyle.update(m.ColorStyle)
	}

	return &style, nil
}

func (s Style) day(t time.Time, cal time.Month) (string, SvgStyle) {
	style := make(SvgStyle)
	switch {
	case s.isHoliday(t):
		style = s.Holiday.style()
	case s.isWeekend(t):
		style = s.Weekend.style()
	case s.isShortday(t):
		style = s.Shortday.style()
	}
	if t.Month() != cal {
		style["opacity"] = s.Ghost
	}
	return strconv.Itoa(t.Day()), style
}

func (s Style) isWeekend(day time.Time) bool {
	if day.Weekday() != time.Saturday && day.Weekday() != time.Sunday {
		return false
	}
	return !s.contains(s.NotWeekend, day)
}

func (s Style) isHoliday(day time.Time) bool {
	return s.contains(s.Holidays, day)
}

func (s Style) isShortday(day time.Time) bool {
	return s.contains(s.Shortdays, day)
}

func (s Style) contains(cfg map[int]map[int][]int, day time.Time) bool {
	if m, found := cfg[day.Year()][int(day.Month())]; found {
		for _, d := range m {
			if d == day.Day() {
				return true
			}
		}
	}
	return false
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
	Weight      string
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

func months() [12]string {
	months := [12]string{}
	for i := 0; i < 12; i++ {
		months[i] = time.Month(i + 1).String()
	}
	return months
}

func replace(s *string, with string) {
	if len(with) > 0 {
		*s = with
	}
}
