package configtx

type Config struct {
	Organizations []Organization `yaml:"Organizations,omitempty"`
	Capabilities  Capabilities   `yaml:"Capabilities,omitempty"`
	Application   Application    `yaml:"Application,omitempty"`
	Orderer       Orderer        `yaml:"Orderer,omitempty"`
	Channel       Channel        `yaml:"Channel,omitempty"`
	Profiles      Profiles       `yaml:"Profiles,omitempty"`
}

type Organization struct {
	Name             string   `yaml:"Name,omitempty"`
	ID               string   `yaml:"ID,omitempty"`
	MSPDir           string   `yaml:"MSPDir,omitempty"`
	Policies         Policies `yaml:"Policies,omitempty"`
	OrdererEndpoints []string `yaml:"OrdererEndpoints,omitempty"`
}

type Policies struct {
	Readers              Policy `yaml:"Readers,omitempty"`
	Writers              Policy `yaml:"Writers,omitempty"`
	Admins               Policy `yaml:"Admins,omitempty"`
	Endorsement          Policy `yaml:"Endorsement,omitempty"`
	LifecycleEndorsement Policy `yaml:"LifecycleEndorsement,omitempty"`
	BlockValidation      Policy `yaml:"BlockValidation,omitempty"`
}

type Policy struct {
	Type string      `yaml:"Type,omitempty"`
	Rule interface{} `yaml:"Rule,omitempty"`
}

type Capabilities struct {
	Channel     Capability `yaml:"Channel,omitempty"`
	Orderer     Capability `yaml:"Orderer,omitempty"`
	Application Capability `yaml:"Application,omitempty"`
}

type Capability struct {
	V2_0 bool `yaml:"V2_0"`
}

type Application struct {
	Organizations []Organization `yaml:"Organizations,omitempty"`
	Policies      Policies       `yaml:"Policies,omitempty"`
	Capabilities  Capability     `yaml:"Capabilities,omitempty"`
}

type Orderer struct {
	OrdererType   string    `yaml:"OrdererType,omitempty"`
	Addresses     []string  `yaml:"Addresses,omitempty"`
	EtcdRaft      EtcdRaft  `yaml:"EtcdRaft,omitempty"`
	BatchTimeout  string    `yaml:"BatchTimeout,omitempty"`
	BatchSize     BatchSize `yaml:"BatchSize,omitempty"`
	Organizations []string  `yaml:"Organizations,omitempty"`
	Policies      Policies  `yaml:"Policies,omitempty"`
}

type EtcdRaft struct {
	Consenters []Consenter `yaml:"Consenters,omitempty"`
}

type Consenter struct {
	Host          string `yaml:"Host,omitempty"`
	Port          int    `yaml:"Port,omitempty"`
	ClientTLSCert string `yaml:"ClientTLSCert,omitempty"`
	ServerTLSCert string `yaml:"ServerTLSCert,omitempty"`
}

type BatchSize struct {
	MaxMessageCount   int    `yaml:"MaxMessageCount,omitempty"`
	AbsoluteMaxBytes  string `yaml:"AbsoluteMaxBytes,omitempty"`
	PreferredMaxBytes string `yaml:"PreferredMaxBytes,omitempty"`
}

type Channel struct {
	Policies     Policies   `yaml:"Policies,omitempty"`
	Capabilities Capability `yaml:"Capabilities,omitempty"`
}

type Profiles struct {
	OrdererGenesis OrdererGenesis `yaml:"OrdererGenesis,omitempty"`
	OrgsChannel    OrgsChannel    `yaml:"OrgsChannel,omitempty"`
}

type OrdererGenesis struct {
	Policies     Policies              `yaml:"Policies,omitempty"`
	Capabilities Capability            `yaml:"Capabilities,omitempty"`
	Orderer      OrdererOrdererGenesis `yaml:"Orderer,omitempty"`
	Consortiums  Consortiums           `yaml:"Consortiums,omitempty"`
}

type Consortiums struct {
	SampleConsortium SampleConsortium `yaml:"SampleConsortium,omitempty"`
}

type SampleConsortium struct {
	Organizations []Organization `yaml:"Organizations,omitempty"`
}

type OrdererOrdererGenesis struct {
	OrdererType   string         `yaml:"OrdererType,omitempty"`
	Addresses     []string       `yaml:"Addresses,omitempty"`
	EtcdRaft      EtcdRaft       `yaml:"EtcdRaft,omitempty"`
	BatchTimeout  string         `yaml:"BatchTimeout,omitempty"`
	BatchSize     BatchSize      `yaml:"BatchSize,omitempty"`
	Policies      Policies       `yaml:"Policies,omitempty"`
	Organizations []Organization `yaml:"Organizations,omitempty"`
	Capabilities  Capability     `yaml:"Capabilities,omitempty"`
}

type OrgsChannel struct {
	Consortium   string      `yaml:"Consortium,omitempty"`
	Policies     Policies    `yaml:"Policies,omitempty"`
	Capabilities Capability  `yaml:"Capabilities,omitempty"`
	Application  Application `yaml:"Application,omitempty"`
}
