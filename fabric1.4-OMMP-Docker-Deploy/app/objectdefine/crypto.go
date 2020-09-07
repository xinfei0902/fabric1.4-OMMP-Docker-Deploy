package objectdefine

type HostnameData struct {
	Prefix string
	Index  int
	Domain string
}

type SpecData struct {
	Hostname   string
	Domain     string
	CommonName string
}
type NodeTemplate struct {
	Count    int      `yaml:"Count,omitempty"`
	Start    int      `yaml:"Start,omitempty"`
	Hostname string   `yaml:"Hostname,omitempty"`
	SANS     []string `yaml:"SANS,omitempty"`
}

type NodeSpec struct {
	Hostname           string   `yaml:"Hostname,omitempty"`
	CommonName         string   `yaml:"CommonName,omitempty"`
	Country            string   `yaml:"Country,omitempty"`
	Province           string   `yaml:"Province,omitempty"`
	Locality           string   `yaml:"Locality,omitempty"`
	OrganizationalUnit string   `yaml:"OrganizationalUnit,omitempty"`
	StreetAddress      string   `yaml:"StreetAddress,omitempty"`
	PostalCode         string   `yaml:"PostalCode,omitempty"`
	SANS               []string `yaml:"SANS,omitempty"`
}

type UsersSpec struct {
	Count int `yaml:"Count"`
}

type OrgSpec struct {
	Name          string       `yaml:"Name"`
	Domain        string       `yaml:"Domain,omitempty"`
	EnableNodeOUs bool         `yaml:"EnableNodeOUs,omitempty"`
	CA            NodeSpec     `yaml:"CA,omitempty"`
	Template      NodeTemplate `yaml:"Template,omitempty"`
	Specs         []NodeSpec   `yaml:"Specs,omitempty"`
	Users         UsersSpec    `yaml:"Users,omitempty"`
}

type CryptoForOrgConfig struct {
	OrdererOrgs []OrgSpec `yaml:"OrdererOrgs"`
	PeerOrgs    []OrgSpec `yaml:"PeerOrgs"`
}
