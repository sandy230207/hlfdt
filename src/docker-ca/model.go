package dockerca

type Config struct {
	Version  string               `yaml:"version,omitempty"`
	Networks Network              `yaml:"networks"`
	Services map[string]Container `yaml:"services,omitempty"`
}

type Container struct {
	Image         string   `yaml:"image,omitempty"`
	Environment   []string `yaml:"environment,omitempty"`
	Ports         []string `yaml:"ports,omitempty"`
	Command       string   `yaml:"command,omitempty"`
	Volumes       []string `yaml:"volumes,omitempty"`
	ContainerName string   `yaml:"container_name,omitempty"`
	Networks      []string `yaml:"networks,omitempty"`
}

type Network struct {
	Test string `yaml:"test"`
}
