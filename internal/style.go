package internal

import (
	"fmt"
	"maps"
	"os"
	"slices"
	"strconv"
	"time"

	"gopkg.in/yaml.v2"
)

type ColorStyle struct {
	Weekday  SvgStyle // Header, like Mn, Tu...
	Number   SvgStyle // Group style of all numbers (days)
	Weekend  SvgStyle
	Holiday  SvgStyle
	Shortday SvgStyle
	Ghost    SvgStyle // Adjacent month days style (opacity)
	Shadow   SvgStyle // Expecting opacity and color of calendar background color
}

func (s *ColorStyle) apply(a ColorStyle) {
	maps.Copy(s.Weekday, a.Weekday)
	maps.Copy(s.Number, a.Number)
	maps.Copy(s.Weekend, a.Weekend)
	maps.Copy(s.Holiday, a.Holiday)
	maps.Copy(s.Shortday, a.Shortday)
	maps.Copy(s.Ghost, a.Ghost)
	maps.Copy(s.Shadow, a.Shadow)
}

type Text struct {
	Pos
	Title string
	Style SvgStyle
}

type Style struct {
	Size   Size
	Extend int // Extend image width and height in mm after draw.
	Grid   Pos
	Pos    Pos

	Layout     LayoutType
	Holidays   []int
	NotWeekend []int
	ShortDays  []int
	Background string

	ColorStyle   ColorStyle
	WeekdayNames [7]string

	Text []Text
}

func DefaultStyle() Style {
	return Style{
		//
		Size: Size{Width: a4Long, Height: a4Short},
		ColorStyle: ColorStyle{
			Weekday:  map[string]string{"fill": "#333333", "font-size": "30", "text-anchor": "middle"},
			Number:   map[string]string{"fill": "#333333", "font-size": "50", "text-anchor": "middle"},
			Weekend:  map[string]string{"fill": "#b03333"},
			Holiday:  map[string]string{"fill": "#b03333"},
			Shortday: map[string]string{"fill": "#584848"},
			Shadow:   map[string]string{"fill": "#ffffff", "opacity": "0.3"},
			Ghost:    map[string]string{"opacity": "0.3"},
		},
		WeekdayNames: weekdays(),
	}
}

func LoadStyle(path string) (*Style, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	style := Style{}
	if err := yaml.UnmarshalStrict(data, &style); err != nil {
		return nil, err
	}

	// Apply defaults
	s := DefaultStyle()
	if style.Size.Width == 0 {
		style.Size.Width = s.Size.Width
	}
	if style.Size.Height == 0 {
		style.Size.Height = s.Size.Height
	}
	if len(style.WeekdayNames[0]) == 0 {
		style.WeekdayNames = s.WeekdayNames
	}
	s.ColorStyle.apply(style.ColorStyle)
	style.ColorStyle = s.ColorStyle

	return &style, nil
}

func (s Style) day(t time.Time, cal time.Month) (string, SvgStyle) {
	style := make(SvgStyle)
	switch {
	case s.isHoliday(t):
		style = maps.Clone(s.ColorStyle.Holiday)
	case s.isWeekend(t):
		style = maps.Clone(s.ColorStyle.Weekend)
	case s.isShortday(t):
		style = maps.Clone(s.ColorStyle.Shortday)
	}
	if t.Month() != cal {
		maps.Copy(style, s.ColorStyle.Ghost)
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
