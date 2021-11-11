package dockercouch

type Config struct {
	Version  string                 `yaml:"version,omitempty"`
	Networks Network                `yaml:"networks"`
	Services map[string]interface{} `yaml:"services,omitempty"`
}

type DBContainer struct {
	Image         string   `yaml:"image,omitempty"`
	Environment   []string `yaml:"environment,omitempty"`
	Ports         []string `yaml:"ports,omitempty"`
	ContainerName string   `yaml:"container_name,omitempty"`
	Networks      []string `yaml:"networks,omitempty"`
}

type Container struct {
	Environment []string `yaml:"environment,omitempty"`
	DependsOn   []string `yaml:"depends_on,omitempty"`
}

type Network struct {
	Test string `yaml:"test"`
}
