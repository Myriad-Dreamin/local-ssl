package main

import (
	"flag"
	"fmt"
	"github.com/Myriad-Dreamin/local-ssl/lib/ssl"
	"os"
	"path/filepath"
)

var commandHelpArgs struct {
	flagSetRef
	usage     func()
	lateUsage func() int
}

func CommandHelp(*ssl.Env) int {
	args := &commandHelpArgs
	if args.flagSet.NArg() == 0 {
		args.usage()
		fmt.Printf("Help Usage:\n  %v help [command] - get usage of `command`\n",
			filepath.Base(os.Args[0]))
		return 0
	}
	return args.lateUsage()
}

func init() {
	fs := flag.NewFlagSet("help", flag.ExitOnError)
	args := &commandHelpArgs
	args.flagSet = fs
	args.usage = usageRoot
	args.lateUsage = func() int {
		var cmd = args.flagSet.Arg(0)
		if command, ok := commands[cmd]; !ok {
			args.usage()
		} else {
			command.set.flagSet.Usage()
		}
		return 0
	}
}
