package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"
)

const (
	VERSION string = "0.0.1"
)

var commands = []cli.Command{
	{
		Name:    "nikkei",
		Aliases: []string{"n"},
		Usage:   "collect nikkei voucher",
		Action:  cmdNikkei,
		Subcommands: []cli.Command{
			{
				Name:   "clean",
				Usage:  "clean nikkei cache directory",
				Action: cmdNikkeiClean,
			},
		},
	},
}

func main() {
	os.Exit(run())
}

func run() int {
	app := cli.NewApp()
	app.Name = "eliminateToil"
	app.Usage = "Erase toil from this world."
	app.Version = VERSION
	app.Commands = commands

	// TODO expalin command help

	// default execute option nikkei
	if len(os.Args) == 1 {
		os.Args = append(os.Args, "nikkei")
	}

	return msg(app.Run(os.Args))
}

/*
func appRun(c *cli.Context) error {
	return nil
}
*/

func msg(err error) int {
	if err != nil {
		errorMsg := fmt.Sprintf("%s: %+v\n", os.Args[0], err)
		fmt.Fprintf(os.Stderr, errorMsg)
		file, _ := os.Create("error.log")
		file.Write([]byte(errorMsg))
		return 1
	}
	return 0
}
