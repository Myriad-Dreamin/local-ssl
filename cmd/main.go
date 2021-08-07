package main

import (
	"fmt"
	"github.com/Myriad-Dreamin/local-ssl/lib/ssl"
	"os"
)

func main() {
	env := ssl.CreateEnv()
	env.CheckBin()
	if env.HasErr() {
		return
	}

	if len(os.Args) < 2 {
		usageRoot()
		panic(fmt.Errorf("expected command positional arguments"))
	}

	if command, ok := commands[os.Args[1]]; !ok {
		panic(fmt.Errorf("invalid command"))
	} else {
		if command.set.flagSet != nil {
			err := command.set.flagSet.Parse(os.Args[2:])
			if err != nil {
				usageRoot()
				panic(err)
			}
		}
		os.Exit(command.entry(env))
	}
}
