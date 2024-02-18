package main

import (
	"os"
	"time"

	"github.com/jessevdk/go-flags"

	"github.com/wolandr/paracal/internal"
)

var opt struct {
	Style  flags.Filename `short:"s" long:"style" description:"Path to style configuration" default:"paracal.yaml"`
	Layout string         `short:"l" long:"layout" choice:"left" choice:"right" choice:"bottom" choice:"top" choice:"square" choice:"square_h" description:"Calendar layout"`
	Back   string         `short:"b" long:"back" description:"Background image path or background color in #hex format"`
	Year   int            `short:"y" long:"year" description:"Year, 0 for current month" default:"0"`
	Month  time.Month     `short:"m" long:"month" description:"Month [1-12], 0 for current month" default:"0"`
	Output flags.Filename `short:"o" long:"output" description:"Output svg file" default:""`
	Ext    int            `short:"e" long:"extend" description:"Extend canvas in mm" default:"0"`
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
		opt.Month = now.Month()
	}

	style, err := internal.LoadStyle(string(opt.Style), opt.Month)
	if err != nil {
		println("Failed to load style configuration,", err.Error())
		os.Exit(1)
	}

	f, err := os.Create(string(opt.Output))
	if err != nil {
		println("Invalid output,", err.Error())
		os.Exit(1)
	}
	defer f.Close()

	internal.Draw(f, opt.Year, opt.Month, internal.LayoutType(opt.Layout), *style, opt.Back, opt.Ext)
}
