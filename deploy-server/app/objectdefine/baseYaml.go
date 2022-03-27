package objectdefine

//BaseYaml  一键部署的结构
type BaseYaml struct {
	Version string                   `yaml:"version"`
	Service map[string]*DockerConfig `yaml:"services"`
}

//DockerConfig 每一个容器配置结构
type DockerConfig struct {
	ContainerName string   `yaml:"container_name"`
	Images        string   `yaml:"image"`
	TTY           bool     `yaml:"tty"`
	StdinOpen     bool     `yaml:"stdin_open"`
	Environment   []string `yaml:"environment"`
	WorkingDir    string   `yaml:"working_dir"`
	Command       string   `yaml:"command"`
	Volumes       []string `yaml:"volumes"`
	Ports         []string `yaml:"ports"`
	DependsOn     []string `yaml:"depends_on"`
	ExtraHosts    []string `yaml:"extra_hosts"`
}
