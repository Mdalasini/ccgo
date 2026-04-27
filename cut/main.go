package main

import (
	"io"
	"os"

	"github.com/alecthomas/kong"
)

type cli struct {
	File      string `arg:"" default:"-" optional:"" help:"Path to the delimited file (default: stdin)."`
	Fields    []int  `short:"f" required:"" help:"Field(s) to extract (1-indexed)."`
	Delimiter string `short:"d" default:"\t" help:"Delimiter character (default: tab)."`
}

func (c *cli) Run() error {
	var r io.Reader
	if c.File == "-" {
		r = os.Stdin
	} else {
		f, err := os.Open(c.File)
		if err != nil {
			return err
		}
		defer f.Close()
		r = f
	}
	return cut(r, os.Stdout, c.Fields, c.Delimiter)
}

func parse() {
	var c cli
	ctx := kong.Parse(&c,
		kong.Name("cut"),
		kong.Description("Cut out selected fields from delimited data."),
	)
	ctx.FatalIfErrorf(ctx.Run())
}

func main() {
	parse()
}
