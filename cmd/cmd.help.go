package main

import (
	"flag"
	"github.com/Myriad-Dreamin/local-ssl/lib/ssl"
)

var commandHelpArgs struct {
	flagSetRef
	usage func()
}

func CommandHelp(*ssl.Env) int {
	commandHelpArgs.usage()
	return 0
}

func init() {
	fs := flag.NewFlagSet("help", flag.ExitOnError)
	args := &commandHelpArgs
	args.flagSet = fs
	args.usage = usageRoot
}
