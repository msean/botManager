package config

type Local struct {
	Path        string `mapstructure:"path" json:"path" yaml:"path"`                         // 本地文件访问路径
	StorePath   string `mapstructure:"store-path" json:"store-path" yaml:"store-path"`       // 本地文件存储路径
	SessionPath string `mapstructure:"session-path" json:"session-path" yaml:"session-path"` // 本地文件存储路径
}
