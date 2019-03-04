package command

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/hcl/hcl/printer"
	"github.com/mitchellh/colorstring"
	"github.com/sergi/go-diff/diffmatchpatch"
)

// FmtCommand is one of the subcommands
type FmtCommand struct {
	CLI
	Option FmtOption
}

// FmtOption is the options for FmtCommand
type FmtOption struct {
	Check bool
	Write bool
}

func (c *FmtCommand) flagSet() *flag.FlagSet {
	flags := flag.NewFlagSet("fmt", flag.ExitOnError)
	flags.BoolVar(&c.Option.Check, "check", false, "perform a syntax check on the given files and produce diagnostics")
	flags.BoolVar(&c.Option.Write, "write", true, "overwrite source files instead of writing to stdout")
	flags.VisitAll(func(f *flag.Flag) {
		if s := os.Getenv(strings.ToUpper(envEnvPrefix + f.Name)); s != "" {
			f.Value.Set(s)
		}
	})
	return flags
}

// Run run fmt command
func (c *FmtCommand) Run(args []string) int {
	flags := c.flagSet()
	if err := flags.Parse(args); err != nil {
		return c.exit(err)
	}

	files := flags.Args()
	return c.exit(c.fmt(files))
}

// Synopsis returns synopsis
func (c *FmtCommand) Synopsis() string {
	return "Formats a policy source to a canonical format."
}

// Help returns help message
func (c *FmtCommand) Help() string {
	var b bytes.Buffer
	flags := c.flagSet()
	flags.SetOutput(&b)
	flags.PrintDefaults()
	return fmt.Sprintf(
		"Usage of %s:\n  %s\n\nOptions:\n%s", flags.Name(), c.Synopsis(), b.String(),
	)
}

func (c *FmtCommand) fmt(files []string) error {
	if len(files) == 0 {
		if c.Option.Write {
			return errors.New("cannot use -w without source files")
		}

		return c.processFile("<stdin>", c.Stdin, c.Stdout)
	}

	for _, file := range files {
		switch dir, err := os.Stat(file); {
		case err != nil:
			return err
		case dir.IsDir():
			// This tool can't walk a whole directory because it doesn't
			// know what file naming schemes will be used by different
			// HCL-embedding applications, so it'll leave that sort of
			// functionality for apps themselves to implement.
			// return fmt.Errorf("can't format directory %s", file)
			c.walkDir(file)
		default:
			if err := c.processFile(file, nil, c.Stdout); err != nil {
				return err
			}
		}
	}
	return nil
}

func (c *FmtCommand) processFile(filename string, in io.Reader, out io.Writer) error {
	if in == nil {
		f, err := os.Open(filename)
		if err != nil {
			return err
		}
		defer f.Close()
		in = f
	}

	src, err := ioutil.ReadAll(in)
	if err != nil {
		return err
	}

	formatted, err := printer.Format(src)
	if err != nil {
		return err
	}

	if c.Option.Check {
		printDiff(lineDiff(string(src), string(formatted)))
		return nil
	}

	if c.Option.Write {
		err = ioutil.WriteFile(filename, formatted, 0644)
	} else {
		_, err = out.Write(formatted)
	}

	return err
}

func isHCL(f os.FileInfo) bool {
	// ignore non-hcl files
	name := f.Name()
	return !f.IsDir() && !strings.HasPrefix(name, ".") && strings.HasSuffix(name, ".hcl")
}

func (c *FmtCommand) walkDir(path string) {
	filepath.Walk(path, c.visitFile)
}

func (c *FmtCommand) visitFile(path string, f os.FileInfo, err error) error {
	if err == nil && isHCL(f) {
		err = c.processFile(path, nil, c.Stdout)
	}

	return err
}

func lineDiff(src1, src2 string) []diffmatchpatch.Diff {
	dmp := diffmatchpatch.New()

	a, b, c := dmp.DiffLinesToChars(src1, src2)
	diffs := dmp.DiffMain(a, b, false)
	result := dmp.DiffCharsToLines(diffs, c)

	return result
}

func printDiff(diffs []diffmatchpatch.Diff) {
	for _, d := range diffs {
		lines := strings.Split(strings.TrimRight(d.Text, "\n"), "\n")
		var prefix string
		switch d.Type {
		case diffmatchpatch.DiffDelete:
			prefix = "[red]- "
		case diffmatchpatch.DiffInsert:
			prefix = "[green]+ "
		case diffmatchpatch.DiffEqual:
			prefix = "  "
		}
		for _, l := range lines {
			s := fmt.Sprintf("%s %s", prefix, l)
			fmt.Println(colorstring.Color(s))
		}
	}
}
