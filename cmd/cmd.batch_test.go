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
apiVersion: certificate.local-ssl.io/v1
scope: skyline-cluster
assets:
  hosts:
  - node2.skyline.io
  - node3.skyline.io
  - node4.skyline.io
mappings:
  ip:
    node2.skyline.io: 192.168.1.5
    node3.skyline.io: 192.168.1.6
    node4.skyline.io: 192.168.1.7
roles:
  ca:
    caConfig:
    keyUsage: [critical, keyCertSign, cRLSign]
  server:
    keyUsage: [critical, extend:critical, serverAuth, nonRepudiation, digitalSignature, keyEncipherment, keyAgreement]
  client:
    keyUsage: [critical, clientAuth, nonRepudiation, digitalSignature, keyEncipherment]
  mTLSClient:
    keyUsage: [critical, clientAuth, nonRepudiation, digitalSignature, keyEncipherment, keyAgreement]
  $default:
    keyUsage: [critical, nonRepudiation, digitalSignature, keyEncipherment]
certs:
  etcd/$hosts:
  - role: $inline
    name: site
  kube-apiserver/$hosts:
  - role: server
`,
	} {
		if code := CommandBatchCreateFromReader(env, bytes.NewReader([]byte(c))); code != 0 {
			t.Errorf("code not zero %d\n", code)
		}
	}
}
