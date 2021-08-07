package main

import (
	"flag"
	"fmt"
	"github.com/Myriad-Dreamin/local-ssl/lib/ssl"
	"os"
	"path/filepath"
	"strings"
)

type flagSetRef struct {
	flagSet *flag.FlagSet
}

var commands = map[string]struct {
	entry func(*ssl.Env) int
	set   *flagSetRef
}{
	"help": {
		entry: CommandHelp,
		set:   &commandHelpArgs.flagSetRef,
	},
}

func usageRoot() {
	var ks = make([]string, 0, len(commands))
	for k := range commands {
		ks = append(ks, k)
	}
	fmt.Printf("Usage: %v [command]\navaliable commands: %s\n",
		filepath.Base(os.Args[0]), strings.Join(ks, " "))
}
