package main

import (
	"fmt"
	"os"
	"time"

	"github.com/jessevdk/go-flags"

	"github.com/wolandr/paracal/internal"
)

var opt struct {
	Style  flags.Filename `short:"s" long:"style" description:"Path to style configuration" default:"paracal.yaml"`
	Layout string         `short:"l" long:"layout" choice:"left" choice:"right" choice:"bottom" choice:"top" choice:"square" choice:"square_v" description:"Calendar layout"`
	Back   string         `short:"b" long:"back" description:"Background image path or background color in #hex format"`
	Year   int            `short:"y" long:"year" description:"Year, 0 for current month" default:"0"`
	Month  int            `short:"m" long:"month" description:"Month [1-12], 0 for current month" base:"10" default:"0"`
	X      int            `long:"posx" description:"X pos" default:"0"`
	Y      int            `long:"posy" description:"Y pos" default:"0"`
	Output flags.Filename `short:"o" long:"output" description:"Output svg file" default:""`
	Debug  bool           `long:"debug" description:"Dump applying style"`
}

func main() {
	parser := flags.NewParser(&opt, flags.Default)
	if _, err := parser.Parse(); err != nil {
		os.Exit(1)
	}

	now := time.Now()
	if opt.Year == 0 {
		opt.Year = now.Year()
	}
	if opt.Month == 0 {
		opt.Month = int(now.Month())
	}

	style, err := internal.LoadStyle(string(opt.Style))
	if err != nil {
		println("Failed to load style configuration,", err.Error())
		os.Exit(1)
	}

	if opt.Year < 100 {
		opt.Year += 2000
	}

	if opt.Layout != "" {
		style.Layout = internal.LayoutType(opt.Layout)
	}

	if opt.Back != "" {
		style.Background = opt.Back
	}

	if opt.X > 0 {
		style.Pos.X = opt.X
	}
	if opt.Y > 0 {
		style.Pos.Y = opt.Y
	}

	f, err := os.Create(string(opt.Output))
	if err != nil {
		println("Invalid output,", err.Error())
		os.Exit(1)
	}
	defer func() { _ = f.Close() }()

	if opt.Debug {
		fmt.Printf("%s\n", style.String())
	}

	internal.Draw(f, opt.Year, time.Month(opt.Month), *style)

	// fmt.Printf("%s\n", opt.Output)
}
