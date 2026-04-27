package main

import (
	"os"

	"github.com/alecthomas/kong"
)

type cli struct {
	File      string `arg:"" help:"Path to the TSV file."`
	Field     int    `short:"f" required:"" help:"Field to extract (1-indexed)."`
	Delimiter string `short:"d" default:"\t" help:"Delimiter character (default: tab)."`
}

func (c *cli) Run() error {
	f, err := os.Open(c.File)
	if err != nil {
		return err
	}
	defer f.Close()
	return cut(f, os.Stdout, c.Field, c.Delimiter)
}

func parse() {
	var c cli
	ctx := kong.Parse(&c,
		kong.Name("cut"),
		kong.Description("Cut out selected fields from delimited data."),
	)
	ctx.FatalIfErrorf(ctx.Run())
}
