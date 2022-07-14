package config

type Config struct {
	Channels               []Channel              `yaml:"Channels"`
	Chaincodes             []Chaincode            `yaml:"Chaincodes"`
	Organizations          []Organization         `yaml:"Organizations"`
	CertificateAuthorities []CertificateAuthority `yaml:"CertificateAuthorities"`
}

type Channel struct {
	Name     string   `yaml:"Name"`
	Policies []Policy `yaml:"Policies"`
}

type Chaincode struct {
	Name     string   `yaml:"Name"`
	Language string   `yaml:"Language"`
	Path     string   `yaml:"Path"`
	Channels []string `yaml:"Channels"`
}

type Organization struct {
	Name string `yaml:"Name"`
	Type string `yaml:"Type"`
	// ConsensusType	string		`yaml:"ConsensusType"`
	BatchTimeout    int       `yaml:"BatchTimeout"`
	BatchSize       BatchSize `yaml:"BatchSize"`
	CommittingPeers []Peer    `yaml:"CommittingPeers"`
	EndorsingPeers  []Peer    `yaml:"EndorsingPeers"`
	Peers           []Peer    `yaml:"Peers"`
	Policies        []Policy  `yaml:"Policies"`
	Channels        []string  `yaml:"Channels"`
}

type BatchSize struct {
	MaxMessageCount   int `yaml:"MaxMessageCount"`
	AbsoluteMaxBytes  int `yaml:"AbsoluteMaxBytes"`
	PreferredMaxBytes int `yaml:"PreferredMaxBytes"`
}

type Peer struct {
	Name    string `yaml:"Name"`
	Address string `yaml:"Address"`
	Port    string `yaml:"Port"`
	DBPort  string `yaml:"DBPort"`
}

type CertificateAuthority struct {
	Name          string   `yaml:"Name"`
	Address       string   `yaml:"Address"`
	Port          string   `yaml:"Port"`
	Organizations []string `yaml:"Organizations"`
}

type Policy struct {
	Name   string `yaml:"Name"`
	Policy string `yaml:"Policy"`
}
