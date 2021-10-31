package main

import (
	"bytes"
	"github.com/Myriad-Dreamin/local-ssl/lib/ssl"
	"testing"
)

func TestCommandBatchCreateFromReader(t *testing.T) {
	env := ssl.CreateEnv()
	env.CheckBin()
	if env.HasErr() {
		return
	}

	for _, c := range []string{
		`
apiVersion: v1
scope: skyline-cluster
roles:
  ca:
    caConfig:
    keyUsage: [critical, keyCertSign, cRLSign]
  server:
    keyUsage: [critical, extend:critical, serverAuth, nonRepudiation, digitalSignature, keyEncipherment, keyAgreement]
  client:
    keyUsage: [critical, extend:critical, clientAuth, nonRepudiation, digitalSignature, keyEncipherment]
  mTLSClient:
    keyUsage: [critical, extend:critical, clientAuth, nonRepudiation, digitalSignature, keyEncipherment, keyAgreement]
`,
	} {
		if code := CommandBatchCreateFromReader(env, bytes.NewReader([]byte(c))); code != 0 {
			t.Errorf("code not zero %d\n", code)
		}
	}
}
