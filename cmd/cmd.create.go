package main

import (
	"flag"
	"fmt"
	"path/filepath"

	"github.com/Myriad-Dreamin/local-ssl/lib/ssl"
)

var commandCreateArgs struct {
	flagSetRef
	projectRoot *string
	site        *string
	org         *string
	unit        *string
	ip          *string
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

	if err := env.SwitchToProject(*args.projectRoot); err != nil {
		panicHelper(err)
	}

	var (
		join        = filepath.Join
		proj        = env.Abs(*args.projectRoot)
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
		sanIP       = *args.ip
		orgName     = *args.org
	)

	if len(sanIP) != 0 {
		sanIP = fmt.Sprintf("IP.1 = %s", sanIP)
	}

	if len(orgName) == 0 {
		orgName = env.CaProjectConfig.O
	}

	env.PushWd(proj)
	env.MakeDir(certs)
	env.MakeDir(siteCerts)
	env.GenerateRSAKey(sitePriLoc)
	env.WriteSignSSLConf(siteConfLoc, &ssl.SignSSLTemplateArgs{
		C:            env.CaProjectConfig.C,
		O:            orgName,
		ST:           env.CaProjectConfig.ST,
		L:            env.CaProjectConfig.L,
		OU:           fmt.Sprintf(`"%s"`, unit),
		CN:           fmt.Sprintf(`"%s"`, site),
		IP:           sanIP,
		EmailAddress: env.CaProjectConfig.EmailAddress,
	})
	env.GenerateCSR(siteConfLoc, sitePriLoc, siteCSRLoc)
	env.CreateSignedCrt(siteConfLoc, siteCSRLoc, siteCrtLoc, caCrtLoc, caPriLoc)
	env.PopWd(proj)

	if env.HasErr() {
		return 1
	}
	return 0
}

func init() {
	fs := flag.NewFlagSet("init", flag.ExitOnError)
	args := &commandCreateArgs
	args.flagSet = fs
	args.projectRoot = fs.String("project", ".", "path to project")
	args.site = fs.String("site", "", "the site that requiring the new certificate")
	args.org = fs.String("org", "", "the unit that requiring the new certificate")
	args.unit = fs.String("unit", "", "the unit that requiring the new certificate")
	args.ip = fs.String("ip", "", "the ip that requiring the new certificate")
}
