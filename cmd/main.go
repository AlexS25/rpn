package main

import (
	"fmt"
	"os"

	"github.com/AlexS25/rpn/internal/application"
)

func usage() {

	fmt.Print(
		"Usage: cmd [arguments]\n\n",
		"Manage calc's list of trusted argumetns\n\n",
		"  --help     - Show current help.\n",
		"  --cmd      - Use command line interface.\n",
		"  --srv      - Run as web-server. (default mod)\n\n",
		"Example of application launch:\n",
		"./cmd --srv\n",
	)

}

type modRun struct {
	isCmd bool
	isSrv bool
}

func (mr *modRun) checkArgs() {

	if len(os.Args) > 1 {
		for _, arg := range os.Args[1:] {
			switch arg {
			case "--help":
				usage()
				os.Exit(0)
			case "--cmd":
				mr.isCmd = true
			case "--srv":
				mr.isSrv = true
			default:
				fmt.Printf("Unknown argument: %q\n", arg)
				os.Exit(1)
			}
		}
	} else {
		mr.isSrv = true
	}

}

func main() {

	var mr modRun
	mr.checkArgs()

	app := application.New()

	if mr.isCmd {
		app.Run()
	} else if mr.isSrv {
		app.RunServer()
	}

}
