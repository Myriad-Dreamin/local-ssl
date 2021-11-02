package main

import (
	"flag"
	"fmt"
	"github.com/Myriad-Dreamin/local-ssl/lib/ssl"
	"io"
	"os"
	"sort"
)

var commandBatchArgs struct {
	flagSetRef
	projectRoot *string
	batchConfig *string
}

type RSAConfig struct {
	Bits uint64 `yaml:"bits"`
}

type CertConfBase struct {
	// keyCertSign (for CA): Subject public key is used to verify signatures on certificates
	// cRLSign (for CA): Subject public key is to verify signatures on revocation information, such as a CRL
	// digitalSignature: Certificate may be used to apply a digital signature
	// nonRepudiation: Certificate may be used to sign data as above,
	//   but the certificate public key may be used to provide non-repudiation services
	// keyEncipherment: Certificate may be used to encrypt a symmetric key which is then transferred to the target
	// dataEncipherment: Certificate may be used to encrypt & decrypt actual application data
	// keyAgreement: Certificate enables use of a key agreement protocol to establish a symmetric key with a target
	// encipherOnly: Public key used only for enciphering data while performing key agreement
	// decipherOnly: Public key used only for deciphering data while performing key agreement
	// extendedKeyUsage:
	//   serverAuth
	//   clientAuth
	//   ipsecEndSystem
	//   ipsecTunnel
	//   ipsecUser
	//   ipsecIKE
	//   codeSigning
	//   emailProtection
	//   timeStamping
	//   OCSPSigning
	//   msCodeInd
	//   msCodeCom
	//   mcCTLSign
	//   msEFS
	KeyUsage  []string   `yaml:"keyUsage"`
	CaConfig  *struct{}  `yaml:"caConfig"`
	Type      string     `yaml:"type"`
	RSAConfig *RSAConfig `yaml:"rsa"`
	O         string     `yaml:"o"`
	CN        string     `yaml:"cn"`
	Sans      []string   `yaml:"sans"`
}

type CertRole struct {
	CertConf CertConfBase `yaml:"roleConf,inline"`
}

type CertConf struct {
	CertConf CertConfBase `yaml:"inlineConf,inline"`
	Role     string       `yaml:"role"`
	Name     string       `yaml:"name"`
}

type BatchCertsConfig struct {
	ApiVersion string                       `yaml:"apiVersion"`
	Scope      string                       `yaml:"scope"`
	Assets     map[string][]string          `yaml:"assets"`
	Mappings   map[string]map[string]string `yaml:"mappings"`
	Roles      map[string]CertRole          `yaml:"roles"`
	Certs      map[string][]CertConf        `yaml:"certs"`
}

func CommandBatchCreateFromReader(env *ssl.Env, r io.Reader) int {
	var conf = getBatchConfig(r)

	if len(conf.Certs) == 0 {
		return 1
	}

	var collected, errors = evaluateBatchConfig(conf)

	sort.Sort(MergedCertConfCmpImpl(collected))
	printEvaluatedConf(conf, collected, false)

	if len(errors) != 0 {
		for _, err := range errors {
			fmt.Printf("error: %e\n", err)
		}
		return 2
	}

	return 0
}

func CommandBatchCreate(env *ssl.Env) int {
	args := &commandBatchArgs
	if len(*args.projectRoot) == 0 {
		args.flagSet.Usage()
		return 2
	}
	if len(*args.batchConfig) == 0 {
		args.flagSet.Usage()
		return 2
	}

	var (
		batchConfig = *args.batchConfig
		code        int
	)

	if err := env.SwitchToProject(*args.projectRoot); err != nil {
		panicHelper(err)
	}

	if batchConfig == "-" {
		code = CommandBatchCreateFromReader(env, os.Stdin)
	} else {
		var f, err = os.OpenFile(batchConfig, os.O_RDONLY, 644)
		panicHelper(err)
		code = CommandBatchCreateFromReader(env, f)
		panicHelper(f.Close())
	}
	return code
}

func init() {
	fs := flag.NewFlagSet("init", flag.ExitOnError)
	args := &commandBatchArgs
	args.flagSet = fs
	args.projectRoot = fs.String("project", ".", "path to project")
	args.batchConfig = fs.String("config", "", "the batch config path")
}
