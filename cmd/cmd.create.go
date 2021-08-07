package main

import (
	"flag"
	"fmt"
	"github.com/Myriad-Dreamin/local-ssl/lib/ssl"
	"path/filepath"
)

var commandCreateArgs struct {
	flagSetRef
	projectRoot *string
	site        *string
	unit        *string
}

func CommandCreate(env *ssl.Env) int {
	if env.HasErr() {
		return 1
	}
	args := &commandCreateArgs
	if len(*args.projectRoot) == 0 {
		args.flagSet.Usage()
		return 2
	}
	if len(*args.site) == 0 {
		args.flagSet.Usage()
		return 2
	}
	if len(*args.unit) == 0 {
		args.flagSet.Usage()
		return 2
	}

	var (
		join        = filepath.Join
		proj        = *args.projectRoot
		site        = *args.site
		unit        = *args.unit
		root        = "."
		caCrtLoc    = join(root, "root-ca.crt")
		caPriLoc    = join(root, "private", "root-ca.key")
		certs       = join(root, "certs")
		siteCerts   = join(certs, site)
		sitePriLoc  = join(siteCerts, "site.key")
		siteConfLoc = join(siteCerts, "site.conf")
		siteCSRLoc  = join(siteCerts, "site.csr")
		siteCrtLoc  = join(siteCerts, "site.crt")
	)

	env.PushWd(proj)
	env.MakeDir(certs)
	env.MakeDir(siteCerts)
	env.GenerateRSAKey(sitePriLoc)
	env.WriteSignSSLConf(siteConfLoc, &ssl.SignSSLTemplateArgs{
		C:            C,
		O:            O,
		ST:           ST,
		L:            L,
		OU:           fmt.Sprintf(`"%s"`, unit),
		CN:           fmt.Sprintf(`"%s"`, site),
		EmailAddress: EmailAddress,
	})
	env.GenerateCSR(siteConfLoc, sitePriLoc, siteCSRLoc)
	env.CreateSignedCrt(siteConfLoc, siteCSRLoc, siteCrtLoc, caCrtLoc, caPriLoc)

	if env.HasErr() {
		return 1
	}
	return 0
}

func init() {
	fs := flag.NewFlagSet("init", flag.ExitOnError)
	args := &commandCreateArgs
	args.flagSet = fs
	args.projectRoot = fs.String("project", "", "path to project")
	args.site = fs.String("site", "", "the site that requiring the new certificate")
	args.unit = fs.String("unit", "", "the unit that requiring the new certificate")
}
