package ssl

import (
	"io"
	"text/template"
)

var signSSLTemplate *template.Template

type SignSSLTemplateArgs struct {
	C            string
	O            string
	ST           string
	L            string
	OU           string
	CN           string
	IP           string
	EmailAddress string
}

func RenderSignSSLConf(w io.Writer, data *SignSSLTemplateArgs) error {
	return signSSLTemplate.Execute(w, data)
}

func init() {
	tmpl := `[default]
name       = root-ca
default_ca = CA_default
name_opt   = ca_default
cert_opt   = ca_default

[CA_default]
home             = .
database         = $home/db/index
serial           = $home/db/serial
crlnumber        = $home/db/crlnumber
certificate      = $home/$name.crt
private_key      = $home/private/$name.key
RANDFILE         = $home/private/random
new_certs_dir    = $home/certs
unique_subject   = no
copy_extensions  = none
default_days     = 365
default_crl_days = 365
default_md       = sha256
policy           = policy_to_match

# Comment out the following two lines for the "traditional"
# (and highly broken) format.
name_opt = ca_default        # Subject Name options
cert_opt = ca_default        # Certificate field options

[policy_to_match]
countryName            = match
stateOrProvinceName    = match
organizationName       = optional
organizationalUnitName = optional
commonName             = supplied
emailAddress           = optional

[req]
default_bits           = 4096
encrypt_key            = no
default_md             = sha256
utf8                   = yes
string_mask            = utf8only
distinguished_name     = req_distinguished_name
prompt                 = no
req_extensions         = v3_req

[v3_req]
basicConstraints = CA:FALSE
keyUsage = critical,keyCertSign,cRLSign
subjectKeyIdentifier = hash
# authorityKeyIdentifier = keyid,issuer
keyUsage = nonRepudiation, digitalSignature, keyEncipherment
subjectAltName = @alt_names

[req_distinguished_name]
C            = {{.C }}
O            = {{.O }}
ST           = {{.ST}}
L            = {{.L }}
OU           = {{.OU}}
CN           = {{.CN}}
emailAddress = {{.EmailAddress}}

[alt_names]
DNS.1        = {{.CN}}
{{.IP}}
`

	signSSLTemplate = template.New("sign-ssl")
	_, err := signSSLTemplate.Parse(tmpl)
	if err != nil {
		panic(err)
	}
}
