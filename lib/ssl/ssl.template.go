package ssl

import (
	"io"
	"text/template"
)

var sslTemplate *template.Template

type SSLTemplateArgs struct {
	C            string
	O            string
	ST           string
	L            string
	OU           string
	CN           string
	EmailAddress string
}

func RenderSSLConf(w io.Writer, data *SSLTemplateArgs) error {
	return sslTemplate.Execute(w, data)
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
organizationName       = match
organizationalUnitName = optional
commonName             = supplied
emailAddress           = optional

[req]
default_bits           = 4096
encrypt_key            = no
default_md             = sha256
utf8                   = yes
string_mask            = utf8only
prompt                 = no
distinguished_name     = distinguished_name
req_extensions         = ca_ext

[distinguished_name]
C            = {{.C }}
O            = {{.O }}
ST           = {{.ST}}
L            = {{.L }}
OU           = {{.OU}}
CN           = {{.CN}}
emailAddress = {{.EmailAddress}}

[ca_ext]
basicConstraints     = critical,CA:true
keyUsage             = critical,keyCertSign,cRLSign
subjectKeyIdentifier = hash
`

	sslTemplate = template.New("ssl")
	_, err := sslTemplate.Parse(tmpl)
	if err != nil {
		panic(err)
	}
}
