package dockernet

type Config struct {
	Version  string               `yaml:"version,omitempty"`
	Volumes  map[string]string    `yaml:"volumes,omitempty"`
	Networks Network              `yaml:"networks"`
	Services map[string]Container `yaml:"services,omitempty"`
}

type Container struct {
	Image         string   `yaml:"image,omitempty"`
	Environment   []string `yaml:"environment,omitempty"`
	Ports         []string `yaml:"ports,omitempty"`
	ContainerName string   `yaml:"container_name,omitempty"`
	Volumes       []string `yaml:"volumes,omitempty"`
	WorkingDir    string   `yaml:"working_dir,omitempty"`
	Command       string   `yaml:"command,omitempty"`
	Networks      []string `yaml:"networks,omitempty"`
	DependsOn     []string `yaml:"depends_on,omitempty"`
	Tty           bool     `yaml:"tty,omitempty"`
	StdinOpen     bool     `yaml:"stdin_open,omitempty"`
}

type Network struct {
	Test string `yaml:"test"`
}
