package main

import (
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/b4b4r07/stein/command"
	"github.com/b4b4r07/stein/pkg/logging"
	"github.com/mitchellh/cli"
)

func main() {
	logWriter, err := logging.LogOutput()
	if err != nil {
		panic(err)
	}
	log.SetOutput(logWriter)

	log.Printf("[INFO] Stein version: %s", command.Version)
	log.Printf("[INFO] Go runtime version: %s", runtime.Version())
	log.Printf("[INFO] CLI args: %#v", os.Args)

	stein := command.CLI{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Meta: &command.Meta{
			UI: &cli.ColoredUi{
				InfoColor:  cli.UiColorBlue,
				ErrorColor: cli.UiColorRed,
				WarnColor:  cli.UiColorYellow,
				Ui: &cli.BasicUi{
					Reader:      os.Stdin,
					Writer:      os.Stdout,
					ErrorWriter: os.Stderr,
				},
			},
		},
	}

	app := &cli.CLI{
		Name:       command.Name,
		Version:    command.Version,
		Args:       os.Args[1:],
		HelpWriter: os.Stdout,
		HelpFunc:   cli.BasicHelpFunc(command.Name),
		Commands: map[string]cli.CommandFactory{
			"apply": func() (cli.Command, error) {
				return &command.ApplyCommand{
					CLI: stein,
				}, nil
			},
			"fmt": func() (cli.Command, error) {
				return &command.FmtCommand{
					CLI: stein,
				}, nil
			},
		},
	}
	exitStatus, err := app.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] %s: %v\n", command.Name, err)
	}
	os.Exit(exitStatus)
}
