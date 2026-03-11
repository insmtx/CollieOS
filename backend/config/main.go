package config

type Config struct {
	Github *GithubAppConfig `yaml:"github,omitempty"`
}
