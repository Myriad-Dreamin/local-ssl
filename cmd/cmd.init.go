package main

import (
	"flag"
	"github.com/Myriad-Dreamin/local-ssl/lib/ssl"
	"path/filepath"
)

var commandInitArgs struct {
	flagSetRef
	projectRoot *string
}

func CommandInit(env *ssl.Env) int {
	if env.HasErr() {
		return 1
	}
	args := &commandInitArgs
	if len(*args.projectRoot) == 0 {
		args.flagSet.Usage()
		return 2
	}

	if err := env.SwitchToProject(*args.projectRoot); err != nil {
		panicHelper(err)
	}

	var (
		join     = filepath.Join
		root     = *args.projectRoot
		confLoc  = join(root, "ssl.conf")
		caLoc    = join(root, "root-ca.crt")
		certs    = join(root, "certs")
		dbDir    = join(root, "db")
		dbIndex  = join(dbDir, "index")
		dbSerial = join(dbDir, "serial")
		priDir   = join(root, "private")
		priLoc   = join(priDir, "root-ca.key")
	)

	env.MakeDir(root)
	env.WriteSSLConf(confLoc, &ssl.SSLTemplateArgs{
		C:            env.CaProjectConfig.C,
		O:            env.CaProjectConfig.O,
		ST:           env.CaProjectConfig.ST,
		L:            env.CaProjectConfig.L,
		OU:           env.CaProjectConfig.CaOU,
		CN:           env.CaProjectConfig.CaCN,
		EmailAddress: env.CaProjectConfig.EmailAddress,
	})

	env.MakeDir(certs)
	env.MakeDir(dbDir)
	env.MakeDir(priDir)
	env.SetFile(dbIndex, "")
	env.SetFile(dbSerial, "01")
	env.GenerateRSAKey(priLoc)
	env.GenerateRootCACrt(confLoc, priLoc, caLoc)

	if env.HasErr() {
		return 1
	}
	return 0
}

func init() {
	fs := flag.NewFlagSet("init", flag.ExitOnError)
	args := &commandInitArgs
	args.flagSet = fs
	args.projectRoot = fs.String("project", ".", "path to project")
}
