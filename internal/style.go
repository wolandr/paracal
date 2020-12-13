package internal

import (
	"io/ioutil"
	"strconv"
	"time"

	"gopkg.in/yaml.v2"
)

type Style struct {
	Weekday    Font
	Number     Font
	Month      Font
	Shadow     Fill
	Holiday    Fill
	Weekend    Fill
	Ghost      string
	Weekdays   [7]string
	Months     [12]string
	Holidays   map[int]map[int][]int
	NotWeekend map[int]map[int][]int
}

func DefaultStyle() Style {
	return Style{
		Weekday:  Font{Size: "60", Color: "#333333"},
		Number:   Font{Size: "80", Weight: "600", Color: "#333333"},
		Month:    Font{Size: "150", Color: "#449955", Stroke: "gray"},
		Shadow:   Fill{Opacity: "0.5", Color: "white"},
		Holiday:  Fill{Color: "#b03333"},
		Weekend:  Fill{Color: "#b03333"},
		Ghost:    "0.4",
		Weekdays: weekdays(),
		Months:   months(),
	}
}

func LoadStyle(path string) (*Style, error) {
	data, err := ioutil.ReadFile(path)
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
		style = s.Holiday.style()
	case s.isWeekend(t):
		style = s.Weekend.style()
	}
	if t.Month() != cal {
		style["opacity"] = s.Ghost
	}
	return strconv.Itoa(t.Day()), style
}

func (s Style) isWeekend(t time.Time) bool {
	if t.Weekday() != time.Saturday && t.Weekday() != time.Sunday {
		return false
	}
	if m, found := s.NotWeekend[t.Year()][int(t.Month())]; found {
		for _, d := range m {
			if d == t.Day() {
				return false
			}
		}
	}
	return true
}

func (s Style) isHoliday(t time.Time) bool {
	if m, found := s.Holidays[t.Year()][int(t.Month())]; found {
		for _, d := range m {
			if d == t.Day() {
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
