package serverconfig

type Config struct {
	Version      string            `yaml:"version,omitempty"`
	Port         int               `yaml:"port,omitempty"`
	Debug        bool              `yaml:"debug,omitempty"`
	CrlsizeLimit int               `yaml:"crlsizelimit,omitempty"`
	TLS          TLS               `yaml:"tls,omitempty"`
	CA           CA                `yaml:"ca,omitempty"`
	CRL          CRL               `yaml:"crl"`
	Registry     Registry          `yaml:"registry"`
	DB           DB                `yaml:"db"`
	LDAP         LDAP              `yaml:"ldap"`
	Affiliations map[string]string `yaml:"affiliations"`
	Signing      Signing           `yaml:"signing"`
	CSR          CSR               `yaml:"csr"`
	BCCSP        BCCSP             `yaml:"bccsp"`
	CACount      string            `yaml:"cacount"`
	CAFiles      string            `yaml:"cafiles"`
	Intermediate Intermediate      `yaml:"intermediate"`
}

type TLS struct {
	Enabled    bool       `yaml:"enabled"`
	Certfile   string     `yaml:"certfile"`
	Keyfile    string     `yaml:"keyfile"`
	ClientAuth ClientAuth `yaml:"clientauth"`
}

type ClientAuth struct {
	Type      string `yaml:"type"`
	Certfiles string `yaml:"certfiles"`
}

type CA struct {
	Name      string `yaml:"name"`
	Keyfile   string `yaml:"keyfile"`
	Certfile  string `yaml:"certfile"`
	Chainfile string `yaml:"chainfile"`
}

type CRL struct {
	Expiry string `yaml:"expiry"`
}

type Registry struct {
	MaxEnrollments int        `yaml:"maxenrollments"`
	Identities     []Identity `yaml:"identities"`
}

type Identity struct {
	Name        string `yaml:"name"`
	Pass        string `yaml:"pass"`
	Type        string `yaml:"Type"`
	Affiliation string `yaml:"affiliation"`
	Attrs       Attrs  `yaml:"attrs"`
}

type Attrs struct {
	HfRegistrarRoles         string `yaml:"hf.Registrar.Roles"`
	HfRegistrarDelegateRoles string `yaml:"hf.Registrar.DelegateRoles"`
	HfRevoker                bool   `yaml:"hf_Revoker"`
	HfIntermediateCA         bool   `yaml:"hf.IntermediateCA"`
	HfGenCRL                 bool   `yaml:"hf.GenCRL"`
	HfRegistrarAttributes    string `yaml:"hf.Registrar.Attributes"`
	HfAffiliationMgr         bool   `yaml:"hf.AffiliationMgr"`
}

type DB struct {
	Type       string `yaml:"type"`
	Datasource string `yaml:"datasource"`
	TLS        TLSDB  `yaml:"tls"`
}

type TLSDB struct {
	Enabled   bool   `yaml:"enabled"`
	Certfiles string `yaml:"certfiles"`
	Client    Client `yaml:"client"`
}

type Client struct {
	Certfile string `yaml:"certfile"`
	Keyfile  string `yaml:"keyfile"`
}

type LDAP struct {
	Enabled   bool      `yaml:"enabled"`
	URL       string    `yaml:"url"`
	TLS       TLSLDAP   `yaml:"tls"`
	Attribute Attribute `yaml:"attribute"`
}

type TLSLDAP struct {
	Certfiles string `yaml:"certfiles"`
	Client    Client `yaml:"client"`
}

type Attribute struct {
	Names      []string         `yaml:"names,flow"`
	Converters []ConverterGroup `yaml:"converters"`
	Maps       Map              `yaml:"maps"`
}

type ConverterGroup struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

type Map struct {
	Groups []ConverterGroup `yaml:"groups"`
}

type Signing struct {
	Default  Default `yaml:"default"`
	Profiles Profile `yaml:"profiles"`
}

type Default struct {
	Usage  []string `yaml:"usage"`
	Expiry string   `yaml:"expiry"`
}

type Profile struct {
	CA  CAProfile `yaml:"ca"`
	TLS Default   `yaml:"tls"`
}

type CAProfile struct {
	Usage        []string     `yaml:"usage"`
	Expiry       string       `yaml:"expiry"`
	CAConstraint Caconstraint `yaml:"caconstraint"`
}

type Caconstraint struct {
	ISCA       bool `yaml:"isca"`
	MaxPathLen int  `yaml:"maxpathlen"`
}

type CSR struct {
	CN    string     `yaml:"cn"`
	Names []NamesCSR `yaml:"names"`
	Hosts []string   `yaml:"hosts"`
	CA    CACSR      `yaml:"ca"`
}

type NamesCSR struct {
	C  string `yaml:"C,omitempty"`
	ST string `yaml:"ST,omitempty"`
	L  string `yaml:"L,omitempty"`
	O  string `yaml:"O,omitempty"`
	OU string `yaml:"OU,omitempty"`
}

type CACSR struct {
	Expiry     string `yaml:"expiry"`
	PathLength int    `yaml:"pathlength"`
}

type BCCSP struct {
	Default string `yaml:"default"`
	SW      SW     `yaml:"sw"`
}

type SW struct {
	Hash         string       `yaml:"hash"`
	Security     int          `yaml:"security"`
	FileKeyStore FileKeyStore `yaml:"filekeystore"`
}

type FileKeyStore struct {
	Keystore string `yaml:"keystore"`
}

type Intermediate struct {
	ParentServer ParentServer `yaml:"parentserver"`
	Enrollment   Enrollment   `yaml:"enrollment"`
	TLS          TLSLDAP      `yaml:"tls"`
}

type ParentServer struct {
	URL    string `yaml:"url"`
	CAName string `yaml:"caname"`
}

type Enrollment struct {
	Hosts   string `yaml:"hosts"`
	Profile string `yaml:"profile"`
	Label   string `yaml:"label"`
}
