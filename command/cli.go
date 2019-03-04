package command

import (
	"io"

	"github.com/mitchellh/cli"
)

const (
	// Name is the application name
	Name = "stein"
	// Version is the application version
	Version = "0.2.4"
)

const (
	envEnvPrefix = "STEIN_"
)

// Meta contains the meta-option that nearly all subcommand inherited
type Meta struct {
	UI cli.Ui
}

// CLI represents the command-line interface
type CLI struct {
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
	Meta   *Meta
}

func (c CLI) exit(msg interface{}) int {
	switch m := msg.(type) {
	case int:
		return m
	case nil:
		return 0
	case string:
		c.Meta.UI.Output(m)
		return 0
	case error:
		c.Meta.UI.Error(m.Error())
		return 1
	default:
		panic(msg)
	}
}
